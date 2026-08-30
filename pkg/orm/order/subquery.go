package order

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// SubQueryOrder models ordering by a scalar subquery.
type SubQueryOrder struct {
	// Query is the sub-query provider yielding the sort value.
	Query contract.QueryStateProvider
	// Direction is the sort direction.
	Direction contract.SortDirection
	// Nulls defines where NULL values are positioned.
	Nulls contract.NullsOrder
}

// GetQuery returns the sub-query provider yielding the scalar sort
// value.
func (o SubQueryOrder) GetQuery() contract.QueryStateProvider {
	return o.Query
}

// GetDirection returns the sorting direction.
func (o SubQueryOrder) GetDirection() contract.SortDirection {
	return o.Direction
}

// GetNulls returns where NULL values should be placed.
func (o SubQueryOrder) GetNulls() contract.NullsOrder {
	return o.Nulls
}
