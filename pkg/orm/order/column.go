package order

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// ColumnOrder models a standard column sort, representing options like
// ORDER BY created_at DESC NULLS LAST.
type ColumnOrder struct {
	// Column is the name of the database column being sorted.
	Column string
	// Direction is the sort direction.
	Direction contract.SortDirection
	// Nulls defines where NULL values are positioned.
	Nulls contract.NullsOrder
	// UsingOperator is the custom comparison operator for sorting.
	UsingOperator contract.ComparisonOperator
}

// GetColumn returns the name of the column being sorted.
func (o ColumnOrder) GetColumn() string {
	return o.Column
}

// GetDirection returns the sorting direction.
func (o ColumnOrder) GetDirection() contract.SortDirection {
	return o.Direction
}

// GetNulls returns where NULL values should be placed.
func (o ColumnOrder) GetNulls() contract.NullsOrder {
	return o.Nulls
}

// GetUsingOperator returns the custom comparison operator used when
// sorting.
func (o ColumnOrder) GetUsingOperator() contract.ComparisonOperator {
	return o.UsingOperator
}
