package contract

// ColumnOrder models a standard column sort (e.g., ORDER BY created_at DESC).
type ColumnOrder interface {
	Order
	// GetColumn returns the name of the column being sorted.
	GetColumn() string
	// GetDirection returns the sorting direction (ASC / DESC).
	GetDirection() SortDirection
	// GetNulls returns where NULL values should be placed.
	GetNulls() NullsOrder
	// GetUsingOperator returns the custom comparison operator when SortUsing
	// is specified.
	GetUsingOperator() ComparisonOperator
}

// RawOrder models sorting by a raw SQL expression (e.g., JSON extraction
// with CAST).
type RawOrder interface {
	Order
	// GetSql returns the raw SQL expression to sort by.
	GetSql() string
	// GetDirection returns the sorting direction.
	GetDirection() SortDirection
}

// SubQueryOrder models ordering by a scalar subquery.
type SubQueryOrder interface {
	Order
	// GetQuery returns the sub-query provider yielding the scalar sort value.
	GetQuery() QueryStateProvider
	// GetDirection returns the sorting direction.
	GetDirection() SortDirection
	// GetNulls returns where NULL values should be placed.
	GetNulls() NullsOrder
}

// Order represents any valid ORDER BY clause.
type Order any

// SortDirection represents the sorting direction or strategy.
type SortDirection string

// NullsOrder controls where NULL values appear in the result set.
type NullsOrder string

// SortAsc represents ascending sort order (ASC).
const SortAsc SortDirection = "ASC"

// SortDesc represents descending sort order (DESC).
const SortDesc SortDirection = "DESC"

// SortUsing represents a custom operator-based sort order (USING %s).
const SortUsing SortDirection = "USING %s"

// NullsFirst specifies that NULL values should sort before non-NULL.
const NullsFirst NullsOrder = "NULLS FIRST"

// NullsLast specifies that NULL values should sort after non-NULL.
const NullsLast NullsOrder = "NULLS LAST"
