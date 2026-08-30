package contract

// Mutator is an interface representing any type used to mutate
// database records.
type Mutator any

// ColumnValue represents a column name and its corresponding value for
// database operations.
type ColumnValue struct {
	// Column is the name of the database column.
	Column string
	// Value is the data value being set or updated in the column.
	Value any
}

// Set represents an optional value used in mutation builders to distinguish
// between an omitted field, a zero value, or NULL.
type Set[T any] struct {
	// IsSet indicates whether the value has been explicitly set.
	IsSet bool
	// Value is the actual generic data value.
	Value T
}
