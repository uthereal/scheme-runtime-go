package multi

// Hydrate maps table columns dynamically into pointers inside model structs.
var Hydrate = struct {
	Role       func(*Role, []string) []any
	Permission func(*Permission, []string) []any
	TenantUser func(*TenantUser, []string) []any
}{
	Role: func(model *Role, columns []string) []any {
		pointers := make([]any, len(columns))
		for i, col := range columns {
			switch col {
			case Schema.Public.Role.ID.ColumnName():
				pointers[i] = &model.ID
			case Schema.Public.Role.Name.ColumnName():
				pointers[i] = &model.Name
			default:
				var discard any
				pointers[i] = &discard
			}
		}
		return pointers
	},
	Permission: func(model *Permission, columns []string) []any {
		pointers := make([]any, len(columns))
		for i, col := range columns {
			switch col {
			case Schema.Public.Permission.ID.ColumnName():
				pointers[i] = &model.ID
			case Schema.Public.Permission.Name.ColumnName():
				pointers[i] = &model.Name
			default:
				var discard any
				pointers[i] = &discard
			}
		}
		return pointers
	},
	TenantUser: func(model *TenantUser, columns []string) []any {
		pointers := make([]any, len(columns))
		for i, col := range columns {
			switch col {
			case Schema.Tenant.TenantUser.ID.ColumnName():
				pointers[i] = &model.ID
			case Schema.Tenant.TenantUser.Email.ColumnName():
				pointers[i] = &model.Email
			case Schema.Tenant.TenantUser.Status.ColumnName():
				pointers[i] = &model.Status
			case Schema.Tenant.TenantUser.RoleID.ColumnName():
				pointers[i] = &model.RoleID
			default:
				var discard any
				pointers[i] = &discard
			}
		}
		return pointers
	},
}
