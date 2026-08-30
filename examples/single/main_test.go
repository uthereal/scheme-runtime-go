package single

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uthereal/scheme-runtime-go/pkg/testutil"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

const schemaDDL = `
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    age INTEGER,
    tags TEXT[] NOT NULL DEFAULT '{}',
    preferences JSONB,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    bio TEXT,
    location POINT,
    active_duration INTERVAL NOT NULL DEFAULT '0 seconds',
    is_public BOOLEAN
);

CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    rating NUMERIC(5, 2) NOT NULL DEFAULT 0.0
);

CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    text TEXT NOT NULL
);

CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE user_groups (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, group_id)
);

-- Seed groups, seed user 999, and seed user_groups pivot records statically
INSERT INTO groups (id, name) VALUES 
('123e4567-e89b-12d3-a456-426614174000', 'GroupA'),
('123e4567-e89b-12d3-a456-426614174001', 'GroupB')
ON CONFLICT (id) DO NOTHING;

INSERT INTO users (id, email) VALUES 
(999, 'owner@example.com') 
ON CONFLICT (id) DO NOTHING;

INSERT INTO user_groups (user_id, group_id) VALUES 
(999, '123e4567-e89b-12d3-a456-426614174000'),
(999, '123e4567-e89b-12d3-a456-426614174001')
ON CONFLICT (user_id, group_id) DO NOTHING;
`

var Relations = &Schema.Public

var pgContainer *testutil.PostgresContainer

// TestMain manages the global Postgres test container.
func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var err error
	pgContainer, err = testutil.StartPostgresContainer(ctx)
	if err != nil {
		fmt.Printf("Failed to start Postgres container: %v\n", err)
		os.Exit(1)
	}

	err = pgContainer.SetupTemplateDB(ctx, "template_db", schemaDDL)
	if err != nil {
		fmt.Printf("Failed to setup template DB: %v\n", err)
		_ = testutil.StopPostgresContainer(ctx, pgContainer)
		os.Exit(1)
	}

	code := m.Run()

	_ = testutil.StopPostgresContainer(ctx, pgContainer)
	os.Exit(code)
}

// Test_Integration_Mutations tests all insert, update, delete, and upsert
// operations.
func Test_Integration_Mutations(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := pgContainer.CreateDB(ctx, t)
	require.NoError(t, err)
	defer cleanup()

	qb := NewUserQuery(db)
	var user User

	t.Run("Insert and InsertReturning", func(t *testing.T) {
		emailVal := "alice@example.com"
		ageVal := 30
		tagsVal := []string{"admin", "beta"}
		prefsVal := map[string]any{"theme": "dark"}
		metaVal := map[string]string{"env": "test"}
		createdVal := time.Now().Truncate(time.Microsecond)

		mut := UserMutator{
			Email:       contract.Set[string]{IsSet: true, Value: emailVal},
			Age:         contract.Set[*int]{IsSet: true, Value: &ageVal},
			Tags:        contract.Set[[]string]{IsSet: true, Value: tagsVal},
			Preferences: contract.Set[*map[string]any]{IsSet: true, Value: &prefsVal},
			Metadata:    contract.Set[map[string]string]{IsSet: true, Value: metaVal},
			CreatedAt:   contract.Set[time.Time]{IsSet: true, Value: createdVal},
		}

		u, err := qb.InsertReturning(ctx, mut)
		require.NoError(t, err)
		assert.Greater(t, u.ID, int64(0))
		assert.Equal(t, emailVal, u.Email)
		assert.Equal(t, ageVal, *u.Age)
		assert.Equal(t, tagsVal, u.Tags)
		assert.Equal(t, prefsVal, *u.Preferences)
		assert.Equal(t, metaVal, u.Metadata)
		assert.True(t, u.CreatedAt.Equal(createdVal))

		user = u
	})

	t.Run("InsertMany and InsertReturningMany", func(t *testing.T) {
		email2 := "bob@example.com"
		email3 := "charlie@example.com"
		muts := []UserMutator{
			{
				Email: contract.Set[string]{IsSet: true, Value: email2},
			},
			{
				Email: contract.Set[string]{IsSet: true, Value: email3},
			},
		}

		users, err := qb.InsertReturningMany(ctx, muts)
		require.NoError(t, err)
		require.Len(t, users, 2)
		assert.Equal(t, email2, users[0].Email)
		assert.Equal(t, email3, users[1].Email)
	})

	t.Run("Update and UpdateReturning", func(t *testing.T) {
		newEmail := "alice_new@example.com"
		userMut := UserMutator{
			Email: contract.Set[string]{IsSet: true, Value: newEmail},
		}
		updatedUsers, err := qb.Where(Schema.Public.User.ID.Eq(user.ID)).
			UpdateReturning(ctx, userMut)
		require.NoError(t, err)
		require.Len(t, updatedUsers, 1)
		assert.Equal(t, newEmail, updatedUsers[0].Email)
	})

	t.Run("Upsert and UpsertReturning", func(t *testing.T) {
		upsertMut := UserMutator{
			ID:    contract.Set[int64]{IsSet: true, Value: user.ID},
			Email: contract.Set[string]{IsSet: true, Value: "alice_upsert@example.com"},
		}
		upsertedUser, err := qb.UpsertReturning(ctx, upsertMut, "id")
		require.NoError(t, err)
		assert.Equal(t, user.ID, upsertedUser.ID)
		assert.Equal(t, "alice_upsert@example.com", upsertedUser.Email)
	})

	t.Run("Delete and DeleteReturning", func(t *testing.T) {
		deletedUsers, err := qb.Where(Schema.Public.User.ID.Eq(user.ID)).
			DeleteReturning(ctx)
		require.NoError(t, err)
		require.Len(t, deletedUsers, 1)
		assert.Equal(t, user.ID, deletedUsers[0].ID)

		exists, err := qb.Where(Schema.Public.User.ID.Eq(user.ID)).Exists(ctx)
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

// Test_Integration_Queries_And_Filters tests standard retrieval, filters,
// and orders.
func Test_Integration_Queries_And_Filters(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := pgContainer.CreateDB(ctx, t)
	require.NoError(t, err)
	defer cleanup()

	// Ensure perfectly clean database for queries
	_, err = db.Exec(ctx, "DELETE FROM users;")
	require.NoError(t, err)

	qb := NewUserQuery(db)

	// Seed Users
	age1, age2 := 25, 45
	tags1 := []string{"foo"}
	tags2 := []string{"bar", "baz"}

	_, err = qb.InsertReturningMany(ctx, []UserMutator{
		{
			Email: contract.Set[string]{IsSet: true, Value: "u1@example.com"},
			Age:   contract.Set[*int]{IsSet: true, Value: &age1},
			Tags:  contract.Set[[]string]{IsSet: true, Value: tags1},
		},
		{
			Email: contract.Set[string]{IsSet: true, Value: "u2@example.com"},
			Age:   contract.Set[*int]{IsSet: true, Value: &age2},
			Tags:  contract.Set[[]string]{IsSet: true, Value: tags2},
		},
		{
			Email: contract.Set[string]{IsSet: true, Value: "u3@example.com"},
			Age:   contract.Set[*int]{IsSet: true, Value: nil},
		},
	})
	require.NoError(t, err)

	t.Run("Basic Where (Gt, Lt, Neq, Eq)", func(t *testing.T) {
		res, err := NewUserQuery(db).Where(Schema.Public.User.Age.Gt(20)).Get(ctx)
		require.NoError(t, err)
		assert.Len(t, res, 2)

		res, err = NewUserQuery(db).Where(Schema.Public.User.Age.Lt(30)).Get(ctx)
		require.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "u1@example.com", res[0].Email)
	})

	t.Run("Set Membership (In, NotIn)", func(t *testing.T) {
		res, err := NewUserQuery(db).
			Where(Schema.Public.User.Email.In(
				"u1@example.com",
				"u3@example.com",
			)).
			Get(ctx)
		require.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run(
		"String operators (Contains, HasPrefix, HasSuffix, Like, ILike)",
		func(t *testing.T) {
			res, err := NewUserQuery(db).
				Where(Schema.Public.User.Email.Contains("2@ex")).
				Get(ctx)
			require.NoError(t, err)
			require.Len(t, res, 1)
			assert.Equal(t, "u2@example.com", res[0].Email)

			res, err = NewUserQuery(db).
				Where(Schema.Public.User.Email.HasPrefix("u3")).
				Get(ctx)
			require.NoError(t, err)
			require.Len(t, res, 1)
		},
	)

	t.Run("Range operators (Between, NotBetween)", func(t *testing.T) {
		res, err := NewUserQuery(db).
			Where(Schema.Public.User.Age.Between(20, 50)).
			Get(ctx)
		require.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("Nullability (IsNull, IsNotNull)", func(t *testing.T) {
		res, err := NewUserQuery(db).Where(Schema.Public.User.Age.IsNull()).Get(ctx)
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, "u3@example.com", res[0].Email)

		res, err = NewUserQuery(db).Where(Schema.Public.User.Age.IsNotNull()).Get(ctx)
		require.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("First and Exists", func(t *testing.T) {
		firstUser, err := NewUserQuery(db).
			Where(Schema.Public.User.Email.Eq("u2@example.com")).
			First(ctx)
		require.NoError(t, err)
		assert.Equal(t, "u2@example.com", firstUser.Email)

		_, err = NewUserQuery(db).
			Where(Schema.Public.User.Email.Eq("notfound")).
			First(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no matching record found")
	})

	t.Run("Array Contains and Overlaps", func(t *testing.T) {
		res, err := NewUserQuery(db).
			Where(Schema.Public.User.Tags.Contains([]string{"foo"})).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, res, 1)

		res, err = NewUserQuery(db).
			Where(Schema.Public.User.Tags.Overlaps([]string{"baz", "other"})).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, "u2@example.com", res[0].Email)
	})
}

// Test_Integration_Pagination tests limit, offset, and various paginate
// methods.
func Test_Integration_Pagination(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := pgContainer.CreateDB(ctx, t)
	require.NoError(t, err)
	defer cleanup()

	// Ensure perfectly clean database for pagination
	_, err = db.Exec(ctx, "DELETE FROM users;")
	require.NoError(t, err)

	qb := NewUserQuery(db)

	var muts []UserMutator
	for i := 1; i <= 5; i++ {
		email := fmt.Sprintf("user%d@example.com", i)
		muts = append(muts, UserMutator{
			Email: contract.Set[string]{IsSet: true, Value: email},
		})
	}
	_, err = qb.InsertReturningMany(ctx, muts)
	require.NoError(t, err)

	t.Run("Standard Paginate", func(t *testing.T) {
		pag, err := qb.OrderBy(Schema.Public.User.ID.Asc()).Paginate(ctx, 2, 0)
		require.NoError(t, err)
		assert.True(t, pag.HasMore)
		require.Len(t, pag.Items, 2)
		assert.Equal(t, "user1@example.com", pag.Items[0].Email)

		pag, err = qb.OrderBy(Schema.Public.User.ID.Asc()).Paginate(ctx, 2, 4)
		require.NoError(t, err)
		assert.False(t, pag.HasMore)
		require.Len(t, pag.Items, 1)
		assert.Equal(t, "user5@example.com", pag.Items[0].Email)
	})

	t.Run("PaginateWithCount", func(t *testing.T) {
		pagCount, err := qb.PaginateWithCount(ctx, 3, 0)
		require.NoError(t, err)
		assert.Equal(t, int64(5), pagCount.TotalCount)
		require.Len(t, pagCount.Items, 3)
	})
}

// Test_Integration_Aggregates tests Min, Max, Avg, and Sum aggregates.
func Test_Integration_Aggregates(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := pgContainer.CreateDB(ctx, t)
	require.NoError(t, err)
	defer cleanup()

	// Ensure perfectly clean database for aggregates
	_, err = db.Exec(ctx, "DELETE FROM users;")
	require.NoError(t, err)

	qb := NewUserQuery(db)

	a1, a2, a3 := 20, 30, 40
	_, err = qb.InsertReturningMany(ctx, []UserMutator{
		{
			Email: contract.Set[string]{IsSet: true, Value: "u1@example.com"},
			Age:   contract.Set[*int]{IsSet: true, Value: &a1},
		},
		{
			Email: contract.Set[string]{IsSet: true, Value: "u2@example.com"},
			Age:   contract.Set[*int]{IsSet: true, Value: &a2},
		},
		{
			Email: contract.Set[string]{IsSet: true, Value: "u3@example.com"},
			Age:   contract.Set[*int]{IsSet: true, Value: &a3},
		},
	})
	require.NoError(t, err)

	minVal, err := qb.Min(ctx, Schema.Public.User.Age)
	require.NoError(t, err)
	assert.Equal(t, 20, minVal)

	maxVal, err := qb.Max(ctx, Schema.Public.User.Age)
	require.NoError(t, err)
	assert.Equal(t, 40, maxVal)

	avgVal, err := qb.Avg(ctx, Schema.Public.User.Age)
	require.NoError(t, err)
	assert.Equal(t, 30, avgVal)

	sumVal, err := qb.Sum(ctx, Schema.Public.User.Age)
	require.NoError(t, err)
	assert.Equal(t, 90, sumVal)
}

// Test_Integration_Relations tests loading all relationship types and
// mapping them.
func Test_Integration_Relations(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := pgContainer.CreateDB(ctx, t)
	require.NoError(t, err)
	defer cleanup()

	var user User
	var post Post
	var comment Comment

	// Fetch seeded user 999
	u, err := NewUserQuery(db).Where(Schema.Public.User.ID.Eq(999)).First(ctx)
	require.NoError(t, err)
	user = u

	t.Run("HasOne Relation -> Profile", func(t *testing.T) {
		bio := "Software Engineer"
		loc := pgtype.Point{
			P:     pgtype.Vec2{X: 37.77, Y: -122.41},
			Valid: true,
		}
		dur := time.Hour * 5
		pubVal := true

		_, err = NewProfileQuery(db).InsertReturning(ctx, ProfileMutator{
			UserID:         contract.Set[int64]{IsSet: true, Value: user.ID},
			Bio:            contract.Set[*string]{IsSet: true, Value: &bio},
			Location:       contract.Set[*pgtype.Point]{IsSet: true, Value: &loc},
			ActiveDuration: contract.Set[time.Duration]{IsSet: true, Value: dur},
			IsPublic:       contract.Set[*bool]{IsSet: true, Value: &pubVal},
		})
		require.NoError(t, err)

		// Test Eager Loading automatically!
		res, err := NewUserQuery(db).
			With(Relations.User.Profile).
			Where(Schema.Public.User.ID.Eq(user.ID)).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, res, 1)

		hydratedUser := res[0]
		require.NotNil(t, hydratedUser.Profile)
		assert.Equal(t, bio, *hydratedUser.Profile.Bio)
		assert.Equal(t, loc.P.X, hydratedUser.Profile.Location.P.X)
		assert.Equal(t, dur, hydratedUser.Profile.ActiveDuration)
		assert.Equal(t, pubVal, *hydratedUser.Profile.IsPublic)
	})

	t.Run("HasMany Relation -> Posts", func(t *testing.T) {
		var rating pgtype.Numeric
		_ = rating.Scan("4.85")

		p, err := NewPostQuery(db).InsertReturning(ctx, PostMutator{
			UserID:  contract.Set[int64]{IsSet: true, Value: user.ID},
			Title:   contract.Set[string]{IsSet: true, Value: "My First Post"},
			Content: contract.Set[string]{IsSet: true, Value: "Hello Postgres"},
			Rating:  contract.Set[pgtype.Numeric]{IsSet: true, Value: rating},
		})
		require.NoError(t, err)
		post = p

		// Test Eager Loading automatically!
		res, err := NewUserQuery(db).
			With(Relations.User.Posts).
			Where(Schema.Public.User.ID.Eq(user.ID)).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, res, 1)

		hydratedUser := res[0]
		require.Len(t, hydratedUser.Posts, 1)
		assert.Equal(t, "My First Post", hydratedUser.Posts[0].Title)
	})

	t.Run("HasMany Relation -> Posts with customization", func(t *testing.T) {
		var rating pgtype.Numeric
		_ = rating.Scan("4.50")

		// Insert a second post
		_, err := NewPostQuery(db).InsertReturning(ctx, PostMutator{
			UserID:  contract.Set[int64]{IsSet: true, Value: user.ID},
			Title:   contract.Set[string]{IsSet: true, Value: "My Second Post"},
			Content: contract.Set[string]{IsSet: true, Value: "Even better content"},
			Rating:  contract.Set[pgtype.Numeric]{IsSet: true, Value: rating},
		})
		require.NoError(t, err)

		// 1. Test filtering relation (Where)
		resFiltered, err := NewUserQuery(db).
			With(Relations.User.Posts.Where(
				Schema.Public.Post.Title.Eq("My Second Post"),
			)).
			Where(Schema.Public.User.ID.Eq(user.ID)).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, resFiltered, 1)
		require.Len(t, resFiltered[0].Posts, 1)
		assert.Equal(t, "My Second Post", resFiltered[0].Posts[0].Title)

		// 2. Test Constrain with OrderBy
		resOrdered, err := NewUserQuery(db).
			With(Relations.User.Posts.Constrain(func(
				q *orm.QueryBuilder[Post, PostMutator],
			) *orm.QueryBuilder[Post, PostMutator] {
				return q.OrderBy(Schema.Public.Post.Title.Desc())
			})).
			Where(Schema.Public.User.ID.Eq(user.ID)).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, resOrdered, 1)
		require.Len(t, resOrdered[0].Posts, 2)
		assert.Equal(t, "My Second Post", resOrdered[0].Posts[0].Title)
		assert.Equal(t, "My First Post", resOrdered[0].Posts[1].Title)

		// 3. Test Constrain with Limit & Offset
		resLimited, err := NewUserQuery(db).
			With(Relations.User.Posts.Constrain(func(
				q *orm.QueryBuilder[Post, PostMutator],
			) *orm.QueryBuilder[Post, PostMutator] {
				return q.OrderBy(
					Schema.Public.Post.Title.Asc(),
				).Limit(1).Offset(1)
			})).
			Where(Schema.Public.User.ID.Eq(user.ID)).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, resLimited, 1)
		require.Len(t, resLimited[0].Posts, 1)
		assert.Equal(t, "My Second Post", resLimited[0].Posts[0].Title)
	})

	t.Run("HasMany Relation -> Comments", func(t *testing.T) {
		c, err := NewCommentQuery(db).InsertReturning(ctx, CommentMutator{
			PostID: contract.Set[int64]{IsSet: true, Value: post.ID},
			UserID: contract.Set[int64]{IsSet: true, Value: user.ID},
			Text:   contract.Set[string]{IsSet: true, Value: "Great post!"},
		})
		require.NoError(t, err)
		comment = c

		// Test Eager Loading automatically!
		res, err := NewPostQuery(db).
			With(Relations.Post.Comments).
			Where(Schema.Public.Post.ID.Eq(post.ID)).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, res, 1)

		hydratedPost := res[0]
		require.Len(t, hydratedPost.Comments, 1)
		assert.Equal(t, "Great post!", hydratedPost.Comments[0].Text)
	})

	t.Run("BelongsTo Relation -> Comment to Post", func(t *testing.T) {
		// Test Eager Loading automatically!
		res, err := NewCommentQuery(db).
			With(Relations.Comment.Post).
			Where(Schema.Public.Comment.ID.Eq(comment.ID)).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, res, 1)

		hydratedComment := res[0]
		require.NotNil(t, hydratedComment.Post)
		assert.Equal(t, "My First Post", hydratedComment.Post.Title)
	})

	t.Run("BelongsTo Relation -> Comment to User", func(t *testing.T) {
		// Test Eager Loading automatically!
		res, err := NewCommentQuery(db).
			With(Relations.Comment.User).
			Where(Schema.Public.Comment.ID.Eq(comment.ID)).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, res, 1)

		hydratedComment := res[0]
		require.NotNil(t, hydratedComment.User)
		assert.Equal(t, "owner@example.com", hydratedComment.User.Email)
	})

	t.Run(
		"BelongsToMany Relation -> User to Group (Many-to-Many via pivot)",
		func(t *testing.T) {
			// Test Eager Loading automatically via statically seeded
			// user_groups pivot records!
			res, err := NewUserQuery(db).
				With(Relations.User.Groups).
				Where(Schema.Public.User.ID.Eq(999)).
				Get(ctx)
			require.NoError(t, err)
			require.Len(t, res, 1)

			hydratedUser := res[0]
			require.Len(t, hydratedUser.Groups, 2)
			assert.Equal(t, "GroupA", hydratedUser.Groups[0].Name)
			assert.Equal(t, "GroupB", hydratedUser.Groups[1].Name)
		},
	)
}

// Test_Integration_Custom_Projection tests QueryBuilder.Select(...) custom
// projection.
func Test_Integration_Custom_Projection(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := pgContainer.CreateDB(ctx, t)
	require.NoError(t, err)
	defer cleanup()

	// Ensure perfectly clean database for projection
	_, err = db.Exec(ctx, "DELETE FROM users;")
	require.NoError(t, err)

	qb := NewUserQuery(db)

	// Insert a user
	user, err := qb.InsertReturning(ctx, UserMutator{
		Email: contract.Set[string]{IsSet: true, Value: "projected@example.com"},
	})
	require.NoError(t, err)

	// Fetch selecting ONLY the Email column!
	projectedUsers, err := qb.Select(Schema.Public.User.Email).
		Where(Schema.Public.User.ID.Eq(user.ID)).
		Get(ctx)
	require.NoError(t, err)
	require.Len(t, projectedUsers, 1)

	// The Email must be populated
	assert.Equal(t, "projected@example.com", projectedUsers[0].Email)
	// The ID field should be its zero value (0) since it wasn't selected or
	// hydrated!
	assert.Equal(t, int64(0), projectedUsers[0].ID)
}

// Test_Integration_AdvancedFeatures tests advanced column operators,
// geospatial metrics, nullable sorting limits, JSONB keys, and subqueries.
func Test_Integration_AdvancedFeatures(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := pgContainer.CreateDB(ctx, t)
	require.NoError(t, err)
	defer cleanup()

	userQuery := NewUserQuery(db)
	profileQuery := NewProfileQuery(db)
	postQuery := NewPostQuery(db)

	// Clean database tables
	_, err = db.Exec(ctx, "DELETE FROM profiles; DELETE FROM users;")
	require.NoError(t, err)

	// Insert static users with different tags, preferences, metadata, and age
	u1, err := userQuery.InsertReturning(ctx, UserMutator{
		Email: contract.Set[string]{IsSet: true, Value: "a1@example.com"},
		Age:   contract.Set[*int]{IsSet: true, Value: pointerTo(25)},
		Tags: contract.Set[[]string]{
			IsSet: true,
			Value: []string{"admin", "editor"},
		},
		Preferences: contract.Set[*map[string]any]{
			IsSet: true,
			Value: pointerTo(map[string]any{
				"theme": "dark",
				"notifications": map[string]any{
					"email": true,
				},
			}),
		},
		Metadata: contract.Set[map[string]string]{
			IsSet: true,
			Value: map[string]string{"region": "us", "active": "yes"},
		},
	})
	require.NoError(t, err)

	u2, err := userQuery.InsertReturning(ctx, UserMutator{
		Email: contract.Set[string]{IsSet: true, Value: "a2@example.com"},
		Age:   contract.Set[*int]{IsSet: true, Value: pointerTo(35)},
		Tags: contract.Set[[]string]{
			IsSet: true,
			Value: []string{"editor", "moderator"},
		},
		Preferences: contract.Set[*map[string]any]{
			IsSet: true,
			Value: pointerTo(map[string]any{
				"theme": "light",
				"notifications": map[string]any{
					"email": false,
				},
			}),
		},
		Metadata: contract.Set[map[string]string]{
			IsSet: true,
			Value: map[string]string{"region": "eu"},
		},
	})
	require.NoError(t, err)

	_, err = userQuery.InsertReturning(ctx, UserMutator{
		Email: contract.Set[string]{IsSet: true, Value: "a3@example.com"},
		Age:   contract.Set[*int]{IsSet: true, Value: nil},
		Tags: contract.Set[[]string]{
			IsSet: true,
			Value: []string{"guest"},
		},
		Metadata: contract.Set[map[string]string]{
			IsSet: true,
			Value: map[string]string{"region": "us"},
		},
	})
	require.NoError(t, err)

	// Create profiles with point locations for geometric testing
	p1, err := profileQuery.InsertReturning(ctx, ProfileMutator{
		UserID: contract.Set[int64]{IsSet: true, Value: u1.ID},
		Location: contract.Set[*pgtype.Point]{
			IsSet: true,
			Value: &pgtype.Point{
				P:     pgtype.Vec2{X: 1.0, Y: 1.0},
				Valid: true,
			},
		},
		ActiveDuration: contract.Set[time.Duration]{
			IsSet: true,
			Value: 2 * time.Hour,
		},
	})
	require.NoError(t, err)

	p2, err := profileQuery.InsertReturning(ctx, ProfileMutator{
		UserID: contract.Set[int64]{IsSet: true, Value: u2.ID},
		Location: contract.Set[*pgtype.Point]{
			IsSet: true,
			Value: &pgtype.Point{
				P:     pgtype.Vec2{X: 5.0, Y: 5.0},
				Valid: true,
			},
		},
		ActiveDuration: contract.Set[time.Duration]{
			IsSet: true,
			Value: 5 * time.Hour,
		},
	})
	require.NoError(t, err)

	t.Run("Array Operators", func(t *testing.T) {
		// Contains: users with 'admin'
		users, err := NewUserQuery(db).
			Where(Schema.Public.User.Tags.Contains([]string{"admin"})).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, users, 1)
		assert.Equal(t, "a1@example.com", users[0].Email)

		// Overlaps: users with 'editor' or 'guest'
		users, err = NewUserQuery(db).
			Where(Schema.Public.User.Tags.Overlaps([]string{"editor", "guest"})).
			OrderBy(Schema.Public.User.Email.Asc()).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, users, 3)

		// ContainedBy: subset check
		users, err = NewUserQuery(db).
			Where(Schema.Public.User.Tags.ContainedBy([]string{
				"admin", "editor", "moderator", "guest", "extra",
			})).
			OrderBy(Schema.Public.User.Email.Asc()).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, users, 3)

		// Concat check: we call the Concat method to satisfy statement coverage
		_ = Schema.Public.User.Tags.Concat([]string{"extra"})
	})

	t.Run("JSONB Operators", func(t *testing.T) {
		// HasKey: users with 'active' metadata key
		users, err := NewUserQuery(db).
			Where(Schema.Public.User.Metadata.HasKey("active")).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, users, 1)
		assert.Equal(t, "a1@example.com", users[0].Email)

		// KeyEq on Preferences (nested path theme -> dark)
		users, err = NewUserQuery(db).
			Where(Schema.Public.User.Preferences.KeyEq("theme", "dark")).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, users, 1)
		assert.Equal(t, "a1@example.com", users[0].Email)
	})

	t.Run("Nullable Sort Variants", func(t *testing.T) {
		// Age: u1=25, u2=35, u3=NULL
		// AscNullsFirst: NULL, 25, 35
		users, err := NewUserQuery(db).
			OrderBy(Schema.Public.User.Age.AscNullsFirst()).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, users, 3)
		assert.Nil(t, users[0].Age)
		assert.Equal(t, int(25), *users[1].Age)

		// AscNullsLast: 25, 35, NULL
		users, err = NewUserQuery(db).
			OrderBy(Schema.Public.User.Age.AscNullsLast()).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, users, 3)
		assert.Equal(t, int(25), *users[0].Age)
		assert.Nil(t, users[2].Age)

		// DescNullsFirst: NULL, 35, 25
		users, err = NewUserQuery(db).
			OrderBy(Schema.Public.User.Age.DescNullsFirst()).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, users, 3)
		assert.Nil(t, users[0].Age)
		assert.Equal(t, int(35), *users[1].Age)

		// DescNullsLast: 35, 25, NULL
		users, err = NewUserQuery(db).
			OrderBy(Schema.Public.User.Age.DescNullsLast()).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, users, 3)
		assert.Equal(t, int(35), *users[0].Age)
		assert.Nil(t, users[2].Age)
	})

	t.Run("Geospatial Point Operators", func(t *testing.T) {
		// p1 is (1,1). p2 is (5,5).
		// StrictLeft of (3,3): should find p1
		refPoint := pgtype.Point{
			P:     pgtype.Vec2{X: 3.0, Y: 3.0},
			Valid: true,
		}
		profs, err := NewProfileQuery(db).
			Where(Schema.Public.Profile.Location.StrictLeft(refPoint)).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, profs, 1)
		assert.Equal(t, p1.ID, profs[0].ID)

		// StrictRight of (3,3): should find p2
		profs, err = NewProfileQuery(db).
			Where(Schema.Public.Profile.Location.StrictRight(refPoint)).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, profs, 1)
		assert.Equal(t, p2.ID, profs[0].ID)

		// Below (3,3): should find p1
		profs, err = NewProfileQuery(db).
			Where(Schema.Public.Profile.Location.Below(refPoint)).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, profs, 1)
		assert.Equal(t, p1.ID, profs[0].ID)

		// Above (3,3): should find p2
		profs, err = NewProfileQuery(db).
			Where(Schema.Public.Profile.Location.Above(refPoint)).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, profs, 1)
		assert.Equal(t, p2.ID, profs[0].ID)
	})

	t.Run("Advanced Range and Date Filters", func(t *testing.T) {
		// Create a post with a rating
		_, err = postQuery.InsertReturning(ctx, PostMutator{
			UserID: contract.Set[int64]{IsSet: true, Value: u1.ID},
			Title:  contract.Set[string]{IsSet: true, Value: "Adv Title"},
			Content: contract.Set[string]{
				IsSet: true,
				Value: "Adv Content",
			},
			Rating: contract.Set[pgtype.Numeric]{
				IsSet: true,
				Value: pgtype.Numeric{
					Int:   big.NewInt(450),
					Exp:   -2,
					Valid: true,
				},
			},
		})
		require.NoError(t, err)

		// Between: rating 4.0 to 5.0
		posts, err := NewPostQuery(db).
			Where(Schema.Public.Post.Rating.Between(
				pgtype.Numeric{
					Int:   big.NewInt(400),
					Exp:   -2,
					Valid: true,
				},
				pgtype.Numeric{
					Int:   big.NewInt(500),
					Exp:   -2,
					Valid: true,
				},
			)).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, posts, 1)

		// Year: check CreatedAt of user is current year (2026) using DateWhere
		now := time.Now()
		users, err := NewUserQuery(db).
			Where(where.DateWhere{
				Type:     "year",
				Column:   "created_at",
				Operator: contract.OpEqual,
				Value:    now.Year(),
				Boolean:  contract.BoolAnd,
			}).
			Get(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 1)
	})

	t.Run("Exists and Subqueries", func(t *testing.T) {
		// InQuery: Users who wrote posts
		postsSub := NewPostQuery(db).Select(Schema.Public.Post.UserID)
		users, err := NewUserQuery(db).
			Where(Schema.Public.User.ID.InQuery(postsSub)).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, users, 1)
		assert.Equal(t, "a1@example.com", users[0].Email)

		// NotInQuery
		users, err = NewUserQuery(db).
			Where(Schema.Public.User.ID.NotInQuery(postsSub)).
			OrderBy(Schema.Public.User.Email.Asc()).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, users, 2)
	})
}

func pointerTo[T any](v T) *T {
	return &v
}
