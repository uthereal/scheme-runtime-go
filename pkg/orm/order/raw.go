package order

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// RawOrder models sorting by a raw SQL expression.
type RawOrder struct {
	// Sql is the raw SQL expression used for ordering.
	Sql string
	// Direction is the sort direction.
	Direction contract.SortDirection
}

// GetSql returns the raw SQL expression to sort by.
func (o RawOrder) GetSql() string {
	return o.Sql
}

// GetDirection returns the sorting direction.
func (o RawOrder) GetDirection() contract.SortDirection {
	return o.Direction
}
