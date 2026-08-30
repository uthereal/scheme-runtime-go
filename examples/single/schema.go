package single

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/column"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/relation"
)

// Schema defines the type-safe column schema metadata for query compiling.
var Schema = struct {
	Public struct {
		User struct {
			ID          column.NumericColumn[User, int64]
			Email       column.StringColumn[User, string]
			Age         column.NullableNumericColumn[User, int]
			Tags        column.ArrayColumn[User, string]
			Preferences column.NullableJSONColumn[User, map[string]any]
			Metadata    column.JSONColumn[User, map[string]string]
			CreatedAt   column.TimestampColumn[User, time.Time]
			Profile     relation.HasOne[User, Profile, ProfileMutator]
			Posts       relation.HasMany[User, Post, PostMutator]
			Groups      relation.BelongsToMany[User, Group, GroupMutator]
		}
		Profile struct {
			ID             column.NumericColumn[Profile, int64]
			UserID         column.NumericColumn[Profile, int64]
			Bio            column.NullableStringColumn[Profile, string]
			Location       column.NullableGeoColumn[Profile, pgtype.Point]
			ActiveDuration column.DurationColumn[Profile, time.Duration]
			IsPublic       column.NullableColumn[Profile, bool]
		}
		Post struct {
			ID       column.NumericColumn[Post, int64]
			UserID   column.NumericColumn[Post, int64]
			Title    column.StringColumn[Post, string]
			Content  column.StringColumn[Post, string]
			Rating   column.DecimalColumn[Post, pgtype.Numeric]
			Comments relation.HasMany[Post, Comment, CommentMutator]
		}
		Comment struct {
			ID     column.NumericColumn[Comment, int64]
			PostID column.NumericColumn[Comment, int64]
			UserID column.NumericColumn[Comment, int64]
			Text   column.StringColumn[Comment, string]
			Post   relation.BelongsTo[Comment, Post, PostMutator]
			User   relation.BelongsTo[Comment, User, UserMutator]
		}
		Group struct {
			ID   column.UUIDColumn[Group]
			Name column.StringColumn[Group, string]
		}
	}
}{
	Public: struct {
		User struct {
			ID          column.NumericColumn[User, int64]
			Email       column.StringColumn[User, string]
			Age         column.NullableNumericColumn[User, int]
			Tags        column.ArrayColumn[User, string]
			Preferences column.NullableJSONColumn[User, map[string]any]
			Metadata    column.JSONColumn[User, map[string]string]
			CreatedAt   column.TimestampColumn[User, time.Time]
			Profile     relation.HasOne[User, Profile, ProfileMutator]
			Posts       relation.HasMany[User, Post, PostMutator]
			Groups      relation.BelongsToMany[User, Group, GroupMutator]
		}
		Profile struct {
			ID             column.NumericColumn[Profile, int64]
			UserID         column.NumericColumn[Profile, int64]
			Bio            column.NullableStringColumn[Profile, string]
			Location       column.NullableGeoColumn[Profile, pgtype.Point]
			ActiveDuration column.DurationColumn[Profile, time.Duration]
			IsPublic       column.NullableColumn[Profile, bool]
		}
		Post struct {
			ID       column.NumericColumn[Post, int64]
			UserID   column.NumericColumn[Post, int64]
			Title    column.StringColumn[Post, string]
			Content  column.StringColumn[Post, string]
			Rating   column.DecimalColumn[Post, pgtype.Numeric]
			Comments relation.HasMany[Post, Comment, CommentMutator]
		}
		Comment struct {
			ID     column.NumericColumn[Comment, int64]
			PostID column.NumericColumn[Comment, int64]
			UserID column.NumericColumn[Comment, int64]
			Text   column.StringColumn[Comment, string]
			Post   relation.BelongsTo[Comment, Post, PostMutator]
			User   relation.BelongsTo[Comment, User, UserMutator]
		}
		Group struct {
			ID   column.UUIDColumn[Group]
			Name column.StringColumn[Group, string]
		}
	}{
		User: struct {
			ID          column.NumericColumn[User, int64]
			Email       column.StringColumn[User, string]
			Age         column.NullableNumericColumn[User, int]
			Tags        column.ArrayColumn[User, string]
			Preferences column.NullableJSONColumn[User, map[string]any]
			Metadata    column.JSONColumn[User, map[string]string]
			CreatedAt   column.TimestampColumn[User, time.Time]
			Profile     relation.HasOne[User, Profile, ProfileMutator]
			Posts       relation.HasMany[User, Post, PostMutator]
			Groups      relation.BelongsToMany[User, Group, GroupMutator]
		}{
			ID:    column.NumericColumn[User, int64]{Name: "id"},
			Email: column.StringColumn[User, string]{Name: "email"},
			Age:   column.NullableNumericColumn[User, int]{Name: "age"},
			Tags:  column.ArrayColumn[User, string]{Name: "tags"},
			Preferences: column.NullableJSONColumn[User, map[string]any]{
				Name: "preferences",
			},
			Metadata:  column.JSONColumn[User, map[string]string]{Name: "metadata"},
			CreatedAt: column.TimestampColumn[User, time.Time]{Name: "created_at"},
		},
		Profile: struct {
			ID             column.NumericColumn[Profile, int64]
			UserID         column.NumericColumn[Profile, int64]
			Bio            column.NullableStringColumn[Profile, string]
			Location       column.NullableGeoColumn[Profile, pgtype.Point]
			ActiveDuration column.DurationColumn[Profile, time.Duration]
			IsPublic       column.NullableColumn[Profile, bool]
		}{
			ID:     column.NumericColumn[Profile, int64]{Name: "id"},
			UserID: column.NumericColumn[Profile, int64]{Name: "user_id"},
			Bio:    column.NullableStringColumn[Profile, string]{Name: "bio"},
			Location: column.NullableGeoColumn[Profile, pgtype.Point]{
				Name: "location",
			},
			ActiveDuration: column.DurationColumn[Profile, time.Duration]{
				Name: "active_duration",
			},
			IsPublic: column.NullableColumn[Profile, bool]{Name: "is_public"},
		},
		Post: struct {
			ID       column.NumericColumn[Post, int64]
			UserID   column.NumericColumn[Post, int64]
			Title    column.StringColumn[Post, string]
			Content  column.StringColumn[Post, string]
			Rating   column.DecimalColumn[Post, pgtype.Numeric]
			Comments relation.HasMany[Post, Comment, CommentMutator]
		}{
			ID:      column.NumericColumn[Post, int64]{Name: "id"},
			UserID:  column.NumericColumn[Post, int64]{Name: "user_id"},
			Title:   column.StringColumn[Post, string]{Name: "title"},
			Content: column.StringColumn[Post, string]{Name: "content"},
			Rating:  column.DecimalColumn[Post, pgtype.Numeric]{Name: "rating"},
		},
		Comment: struct {
			ID     column.NumericColumn[Comment, int64]
			PostID column.NumericColumn[Comment, int64]
			UserID column.NumericColumn[Comment, int64]
			Text   column.StringColumn[Comment, string]
			Post   relation.BelongsTo[Comment, Post, PostMutator]
			User   relation.BelongsTo[Comment, User, UserMutator]
		}{
			ID:     column.NumericColumn[Comment, int64]{Name: "id"},
			PostID: column.NumericColumn[Comment, int64]{Name: "post_id"},
			UserID: column.NumericColumn[Comment, int64]{Name: "user_id"},
			Text:   column.StringColumn[Comment, string]{Name: "text"},
		},
		Group: struct {
			ID   column.UUIDColumn[Group]
			Name column.StringColumn[Group, string]
		}{
			ID:   column.UUIDColumn[Group]{Name: "id"},
			Name: column.StringColumn[Group, string]{Name: "name"},
		},
	},
}
