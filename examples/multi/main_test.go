package multi

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uthereal/scheme-runtime-go/pkg/testutil"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

const multiSchemaDDL = `
CREATE SCHEMA IF NOT EXISTS tenant;

CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE permissions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE tenant.users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE
);

CREATE TABLE tenant.user_permissions (
    user_id BIGINT NOT NULL REFERENCES tenant.users(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);
`

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

	err = pgContainer.SetupTemplateDB(ctx, "multi_template_db", multiSchemaDDL)
	if err != nil {
		fmt.Printf("Failed to setup template DB: %v\n", err)
		_ = testutil.StopPostgresContainer(ctx, pgContainer)
		os.Exit(1)
	}

	code := m.Run()

	_ = testutil.StopPostgresContainer(ctx, pgContainer)
	os.Exit(code)
}

// Test_Integration_MultiSchema verifies cross-schema insertions, queries,
// custom enums, and bidirectional cross-schema eager loading.
func Test_Integration_MultiSchema(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := pgContainer.CreateDB(ctx, t)
	require.NoError(t, err)
	defer cleanup()

	roleQuery := NewRoleQuery(db)
	permQuery := NewPermissionQuery(db)
	userQuery := NewTenantUserQuery(db)

	var adminRole Role
	var guestRole Role
	var readPerm Permission
	var writePerm Permission

	t.Run("Insert roles, permissions and tenant users", func(t *testing.T) {
		r1, err := roleQuery.InsertReturning(ctx, RoleMutator{
			Name: contract.Set[string]{IsSet: true, Value: "Administrator"},
		})
		require.NoError(t, err)
		assert.Greater(t, r1.ID, int64(0))
		assert.Equal(t, "Administrator", r1.Name)
		adminRole = r1

		r2, err := roleQuery.InsertReturning(ctx, RoleMutator{
			Name: contract.Set[string]{IsSet: true, Value: "Guest"},
		})
		require.NoError(t, err)
		assert.Greater(t, r2.ID, int64(0))
		assert.Equal(t, "Guest", r2.Name)
		guestRole = r2

		p1, err := permQuery.InsertReturning(ctx, PermissionMutator{
			Name: contract.Set[string]{IsSet: true, Value: "read"},
		})
		require.NoError(t, err)
		assert.Greater(t, p1.ID, int64(0))
		assert.Equal(t, "read", p1.Name)
		readPerm = p1

		p2, err := permQuery.InsertReturning(ctx, PermissionMutator{
			Name: contract.Set[string]{IsSet: true, Value: "write"},
		})
		require.NoError(t, err)
		assert.Greater(t, p2.ID, int64(0))
		assert.Equal(t, "write", p2.Name)
		writePerm = p2

		u1, err := userQuery.InsertReturning(ctx, TenantUserMutator{
			Email: contract.Set[string]{
				IsSet: true,
				Value: "user1@tenant.com",
			},
			Status: contract.Set[UserStatus]{
				IsSet: true,
				Value: UserStatusActive,
			},
			RoleID: contract.Set[int64]{
				IsSet: true,
				Value: adminRole.ID,
			},
		})
		require.NoError(t, err)
		assert.Greater(t, u1.ID, int64(0))

		u2, err := userQuery.InsertReturning(ctx, TenantUserMutator{
			Email: contract.Set[string]{
				IsSet: true,
				Value: "user2@tenant.com",
			},
			Status: contract.Set[UserStatus]{
				IsSet: true,
				Value: UserStatusInactive,
			},
			RoleID: contract.Set[int64]{
				IsSet: true,
				Value: adminRole.ID,
			},
		})
		require.NoError(t, err)
		assert.Greater(t, u2.ID, int64(0))

		u3, err := userQuery.InsertReturning(ctx, TenantUserMutator{
			Email: contract.Set[string]{
				IsSet: true,
				Value: "user3@tenant.com",
			},
			Status: contract.Set[UserStatus]{
				IsSet: true,
				Value: UserStatusActive,
			},
			RoleID: contract.Set[int64]{
				IsSet: true,
				Value: guestRole.ID,
			},
		})
		require.NoError(t, err)
		assert.Greater(t, u3.ID, int64(0))

		// Insert tenant.user_permissions pivot entries
		_, err = db.Exec(
			ctx,
			"INSERT INTO tenant.user_permissions (user_id, permission_id) VALUES ($1, $2), ($3, $4), ($5, $6)",
			u1.ID, readPerm.ID,
			u1.ID, writePerm.ID,
			u3.ID, readPerm.ID,
		)
		require.NoError(t, err)
	})

	t.Run("Eager Load: TenantUser -> Role (BelongsTo)", func(t *testing.T) {
		users, err := NewTenantUserQuery(db).
			With(Schema.Tenant.TenantUser.Role).
			Where(Schema.Tenant.TenantUser.Email.Eq("user1@tenant.com")).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, users, 1)

		u := users[0]
		require.NotNil(t, u.Role)
		assert.Equal(t, adminRole.ID, u.Role.ID)
		assert.Equal(t, "Administrator", u.Role.Name)
	})

	t.Run("Eager Load: Role -> TenantUser (HasMany)", func(t *testing.T) {
		roles, err := NewRoleQuery(db).
			With(Schema.Public.Role.Users).
			Where(Schema.Public.Role.ID.Eq(adminRole.ID)).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, roles, 1)

		r := roles[0]
		require.Len(t, r.Users, 2)
		assert.Equal(t, "user1@tenant.com", r.Users[0].Email)
		assert.Equal(t, "user2@tenant.com", r.Users[1].Email)
	})

	t.Run(
		"Eager Load: TenantUser -> Permissions (BelongsToMany)",
		func(t *testing.T) {
		users, err := NewTenantUserQuery(db).
			With(Schema.Tenant.TenantUser.Permissions).
			Where(Schema.Tenant.TenantUser.Email.Eq("user1@tenant.com")).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, users, 1)

		u := users[0]
		require.Len(t, u.Permissions, 2)
		assert.Equal(t, "read", u.Permissions[0].Name)
		assert.Equal(t, "write", u.Permissions[1].Name)

		users3, err := NewTenantUserQuery(db).
			With(Schema.Tenant.TenantUser.Permissions).
			Where(Schema.Tenant.TenantUser.Email.Eq("user3@tenant.com")).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, users3, 1)

		u3 := users3[0]
		require.Len(t, u3.Permissions, 1)
		assert.Equal(t, "read", u3.Permissions[0].Name)
	})

	t.Run("Cross-Schema Enum Filtering", func(t *testing.T) {
		activeUsers, err := NewTenantUserQuery(db).
			Where(Schema.Tenant.TenantUser.Status.Eq(UserStatusActive)).
			OrderBy(Schema.Tenant.TenantUser.Email.Asc()).
			Get(ctx)
		require.NoError(t, err)
		require.Len(t, activeUsers, 2)
		assert.Equal(t, "user1@tenant.com", activeUsers[0].Email)
		assert.Equal(t, "user3@tenant.com", activeUsers[1].Email)
	})

	t.Run("Cross-Schema Pagination", func(t *testing.T) {
		pag, err := NewTenantUserQuery(db).
			OrderBy(Schema.Tenant.TenantUser.Email.Asc()).
			Paginate(ctx, 1, 1)
		require.NoError(t, err)
		assert.True(t, pag.HasMore)
		require.Len(t, pag.Items, 1)
		assert.Equal(t, "user2@tenant.com", pag.Items[0].Email)
	})
}
