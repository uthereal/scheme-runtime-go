package multi

// Role represents a role inside the public.roles table.
type Role struct {
	ID    int64
	Name  string
	Users []*TenantUser
}

// Permission represents a shared security permission in public.permissions.
type Permission struct {
	ID   int64
	Name string
}

// TenantUser represents a tenant user inside the tenant.users table.
type TenantUser struct {
	ID          int64
	Email       string
	Status      UserStatus
	RoleID      int64
	Role        *Role
	Permissions []*Permission
}
