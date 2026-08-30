package multi

// Dehydrate maps mutator payloads to column slices and parameter slices.
var Dehydrate = struct {
	Role       func(*RoleMutator) ([]string, []any)
	Permission func(*PermissionMutator) ([]string, []any)
	TenantUser func(*TenantUserMutator) ([]string, []any)
}{
	Role: func(mutator *RoleMutator) ([]string, []any) {
		var cols []string
		var vals []any
		if mutator.ID.IsSet {
			cols = append(cols, Schema.Public.Role.ID.ColumnName())
			vals = append(vals, mutator.ID.Value)
		}
		if mutator.Name.IsSet {
			cols = append(cols, Schema.Public.Role.Name.ColumnName())
			vals = append(vals, mutator.Name.Value)
		}
		return cols, vals
	},
	Permission: func(mutator *PermissionMutator) ([]string, []any) {
		var cols []string
		var vals []any
		if mutator.ID.IsSet {
			cols = append(cols, Schema.Public.Permission.ID.ColumnName())
			vals = append(vals, mutator.ID.Value)
		}
		if mutator.Name.IsSet {
			cols = append(cols, Schema.Public.Permission.Name.ColumnName())
			vals = append(vals, mutator.Name.Value)
		}
		return cols, vals
	},
	TenantUser: func(mutator *TenantUserMutator) ([]string, []any) {
		var cols []string
		var vals []any
		if mutator.ID.IsSet {
			cols = append(cols, Schema.Tenant.TenantUser.ID.ColumnName())
			vals = append(vals, mutator.ID.Value)
		}
		if mutator.Email.IsSet {
			cols = append(
				cols,
				Schema.Tenant.TenantUser.Email.ColumnName(),
			)
			vals = append(vals, mutator.Email.Value)
		}
		if mutator.Status.IsSet {
			cols = append(
				cols,
				Schema.Tenant.TenantUser.Status.ColumnName(),
			)
			vals = append(vals, mutator.Status.Value)
		}
		if mutator.RoleID.IsSet {
			cols = append(
				cols,
				Schema.Tenant.TenantUser.RoleID.ColumnName(),
			)
			vals = append(vals, mutator.RoleID.Value)
		}
		return cols, vals
	},
}
