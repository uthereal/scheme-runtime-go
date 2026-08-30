package multi

import (
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// RoleMutator manages mutations for Role fields.
type RoleMutator struct {
	ID   contract.Set[int64]
	Name contract.Set[string]
}

// PermissionMutator manages mutations for Permission fields.
type PermissionMutator struct {
	ID   contract.Set[int64]
	Name contract.Set[string]
}

// TenantUserMutator manages mutations for TenantUser fields.
type TenantUserMutator struct {
	ID     contract.Set[int64]
	Email  contract.Set[string]
	Status contract.Set[UserStatus]
	RoleID contract.Set[int64]
}
