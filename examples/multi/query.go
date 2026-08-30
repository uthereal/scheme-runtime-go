package multi

import (
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/grammar"
	"github.com/uthereal/scheme-runtime-go/pkg/orm"
)

// NewRoleQuery constructs a QueryBuilder for Role.
func NewRoleQuery(db contract.DB) *orm.QueryBuilder[Role, RoleMutator] {
	compiler := grammar.NewPostgresGrammar()
	return orm.NewQueryBuilder[Role, RoleMutator](
		db,
		compiler,
		Table.Public.Role,
		Hydrate.Role,
		Dehydrate.Role,
	)
}

// NewPermissionQuery constructs a QueryBuilder for Permission.
func NewPermissionQuery(
	db contract.DB,
) *orm.QueryBuilder[Permission, PermissionMutator] {
	compiler := grammar.NewPostgresGrammar()
	return orm.NewQueryBuilder[Permission, PermissionMutator](
		db,
		compiler,
		Table.Public.Permission,
		Hydrate.Permission,
		Dehydrate.Permission,
	)
}

// NewTenantUserQuery constructs a QueryBuilder for TenantUser.
func NewTenantUserQuery(
	db contract.DB,
) *orm.QueryBuilder[TenantUser, TenantUserMutator] {
	compiler := grammar.NewPostgresGrammar()
	return orm.NewQueryBuilder[TenantUser, TenantUserMutator](
		db,
		compiler,
		Table.Tenant.TenantUser,
		Hydrate.TenantUser,
		Dehydrate.TenantUser,
	)
}
