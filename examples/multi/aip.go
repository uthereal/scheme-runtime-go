package multi

import (
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/aip160/fields"
)

// GlobalTenantUserAip160Fields defines global search fields for TenantUser.
var GlobalTenantUserAip160Fields = []contract.Aip160Field[contract.Where]{
	fields.NewStringField[
		string,
		contract.Where,
	](Schema.Tenant.TenantUser.Email),
}

// MapTenantUserColumnToAip160Field maps TenantUser columns to AIP-160
// fields.
var MapTenantUserColumnToAip160Field = map[string]contract.Aip160Field[
	contract.Where,
]{
	"id": fields.NewIntegerField[
		int64,
		contract.Where,
	](Schema.Tenant.TenantUser.ID),
	"email": fields.NewStringField[
		string,
		contract.Where,
	](Schema.Tenant.TenantUser.Email),
	"status": fields.NewStringEnumField[
		UserStatus,
		contract.Where,
	](Schema.Tenant.TenantUser.Status),
}
