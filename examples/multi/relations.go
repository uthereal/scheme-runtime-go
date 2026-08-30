package multi

import (
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/relation"
)

func init() {
	Schema.Tenant.TenantUser.Role = relation.BelongsTo[
		TenantUser,
		Role,
		RoleMutator,
	]{
		Relation: relation.Relation[TenantUser, Role, RoleMutator]{
			ForeignKeyColumns: []contract.Column[Role]{
				Schema.Public.Role.ID,
			},
			ChildQueryFactory: func(
				db contract.DB,
			) *orm.QueryBuilder[Role, RoleMutator] {
				return NewRoleQuery(db)
			},
			LocalKeyExtractor: func(c *TenantUser) any {
				return c.RoleID
			},
			ForeignKeyExtractor: func(p *Role) any {
				return p.ID
			},
			Hydrator: func(c *TenantUser, children []Role) {
				if len(children) > 0 {
					c.Role = &children[0]
				}
			},
		},
	}

	Schema.Public.Role.Users = relation.HasMany[
		Role,
		TenantUser,
		TenantUserMutator,
	]{
		Relation: relation.Relation[Role, TenantUser, TenantUserMutator]{
			ForeignKeyColumns: []contract.Column[TenantUser]{
				Schema.Tenant.TenantUser.RoleID,
			},
			ChildQueryFactory: func(
				db contract.DB,
			) *orm.QueryBuilder[TenantUser, TenantUserMutator] {
				return NewTenantUserQuery(db)
			},
			LocalKeyExtractor: func(p *Role) any {
				return p.ID
			},
			ForeignKeyExtractor: func(c *TenantUser) any {
				return c.RoleID
			},
			Hydrator: func(p *Role, children []TenantUser) {
				ptrs := make([]*TenantUser, len(children))
				for i := range children {
					ptrs[i] = &children[i]
				}
				p.Users = ptrs
			},
		},
	}

	Schema.Tenant.TenantUser.Permissions = relation.BelongsToMany[
		TenantUser,
		Permission,
		PermissionMutator,
	]{
		PivotSchema:             "tenant",
		PivotTable:              "user_permissions",
		PivotForeignKeyToParent: "user_id",
		PivotForeignKeyToChild:  "permission_id",
		ChildForeignKeyColumn:   Schema.Public.Permission.ID,
		ChildQueryFactory: func(
			db contract.DB,
		) *orm.QueryBuilder[Permission, PermissionMutator] {
			return NewPermissionQuery(db)
		},
		LocalKeyExtractor: func(p *TenantUser) any {
			return p.ID
		},
		ChildKeyExtractor: func(c *Permission) any {
			return c.ID
		},
		Hydrator: func(p *TenantUser, children []Permission) {
			ptrs := make([]*Permission, len(children))
			for i := range children {
				ptrs[i] = &children[i]
			}
			p.Permissions = ptrs
		},
	}
}
