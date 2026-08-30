package multi

import (
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Table defines table metadata.
var Table = struct {
	Public struct {
		Role       contract.TableMetadata[Role]
		Permission contract.TableMetadata[Permission]
	}
	Tenant struct {
		TenantUser contract.TableMetadata[TenantUser]
	}
}{
	Public: struct {
		Role       contract.TableMetadata[Role]
		Permission contract.TableMetadata[Permission]
	}{
		Role: contract.TableMetadata[Role]{
			SchemaName: "public",
			TableName:  "roles",
			DefaultColumns: []contract.Column[Role]{
				Schema.Public.Role.ID,
				Schema.Public.Role.Name,
			},
		},
		Permission: contract.TableMetadata[Permission]{
			SchemaName: "public",
			TableName:  "permissions",
			DefaultColumns: []contract.Column[Permission]{
				Schema.Public.Permission.ID,
				Schema.Public.Permission.Name,
			},
		},
	},
	Tenant: struct {
		TenantUser contract.TableMetadata[TenantUser]
	}{
		TenantUser: contract.TableMetadata[TenantUser]{
			SchemaName: "tenant",
			TableName:  "users",
			DefaultColumns: []contract.Column[TenantUser]{
				Schema.Tenant.TenantUser.ID,
				Schema.Tenant.TenantUser.Email,
				Schema.Tenant.TenantUser.Status,
				Schema.Tenant.TenantUser.RoleID,
			},
		},
	},
}
