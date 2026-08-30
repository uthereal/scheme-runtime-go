package multi

import (
	"github.com/uthereal/scheme-runtime-go/pkg/orm/column"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/relation"
)

// Schema defines the type-safe column schema metadata.
var Schema = struct {
	Public struct {
		Role struct {
			ID    column.NumericColumn[Role, int64]
			Name  column.StringColumn[Role, string]
			Users relation.HasMany[Role, TenantUser, TenantUserMutator]
		}
		Permission struct {
			ID   column.NumericColumn[Permission, int64]
			Name column.StringColumn[Permission, string]
		}
	}
	Tenant struct {
		TenantUser struct {
			ID          column.NumericColumn[TenantUser, int64]
			Email       column.StringColumn[TenantUser, string]
			Status      column.StringColumn[TenantUser, UserStatus]
			RoleID      column.NumericColumn[TenantUser, int64]
			Role        relation.BelongsTo[TenantUser, Role, RoleMutator]
			Permissions relation.BelongsToMany[
				TenantUser,
				Permission,
				PermissionMutator,
			]
		}
	}
}{
	Public: struct {
		Role struct {
			ID    column.NumericColumn[Role, int64]
			Name  column.StringColumn[Role, string]
			Users relation.HasMany[Role, TenantUser, TenantUserMutator]
		}
		Permission struct {
			ID   column.NumericColumn[Permission, int64]
			Name column.StringColumn[Permission, string]
		}
	}{
		Role: struct {
			ID    column.NumericColumn[Role, int64]
			Name  column.StringColumn[Role, string]
			Users relation.HasMany[Role, TenantUser, TenantUserMutator]
		}{
			ID: column.NumericColumn[Role, int64]{
				Name: "id",
			},
			Name: column.StringColumn[Role, string]{
				Name: "name",
			},
		},
		Permission: struct {
			ID   column.NumericColumn[Permission, int64]
			Name column.StringColumn[Permission, string]
		}{
			ID: column.NumericColumn[Permission, int64]{
				Name: "id",
			},
			Name: column.StringColumn[Permission, string]{
				Name: "name",
			},
		},
	},
	Tenant: struct {
		TenantUser struct {
			ID          column.NumericColumn[TenantUser, int64]
			Email       column.StringColumn[TenantUser, string]
			Status      column.StringColumn[TenantUser, UserStatus]
			RoleID      column.NumericColumn[TenantUser, int64]
			Role        relation.BelongsTo[TenantUser, Role, RoleMutator]
			Permissions relation.BelongsToMany[
				TenantUser,
				Permission,
				PermissionMutator,
			]
		}
	}{
		TenantUser: struct {
			ID          column.NumericColumn[TenantUser, int64]
			Email       column.StringColumn[TenantUser, string]
			Status      column.StringColumn[TenantUser, UserStatus]
			RoleID      column.NumericColumn[TenantUser, int64]
			Role        relation.BelongsTo[TenantUser, Role, RoleMutator]
			Permissions relation.BelongsToMany[
				TenantUser,
				Permission,
				PermissionMutator,
			]
		}{
			ID: column.NumericColumn[TenantUser, int64]{
				Name: "id",
			},
			Email: column.StringColumn[TenantUser, string]{
				Name: "email",
			},
			Status: column.StringColumn[TenantUser, UserStatus]{
				Name: "status",
			},
			RoleID: column.NumericColumn[TenantUser, int64]{
				Name: "role_id",
			},
		},
	},
}
