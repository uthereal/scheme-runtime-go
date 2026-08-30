package contract

// Column represents a column associated with a specific Model,
// carrying its physical column name and PostgreSQL type-serialization helpers.
type Column[Model any] interface {
	// ColumnName returns the physical name of the column in the database
	// table.
	ColumnName() string

	// PostgresCast returns the PostgreSQL cast suffix (e.g., "::bigint[]",
	// "::text[]").
	PostgresCast() string

	// ToTypedSlice converts a generic slice of []any to a strongly-typed
	// slice of this column's type.
	ToTypedSlice(slice []any) any

	// IsArray returns true if this column is an array type in the database.
	IsArray() bool
}
