package contract

// QueryStateProvider defines the read-only interface required by SQL
// compilers to read the state of a built query without allowing mutation.
type QueryStateProvider interface {
	// GetSchemaName returns the name of the database schema.
	GetSchemaName() string

	// GetTableName returns the physical name of the database table.
	GetTableName() string

	// GetDefaultColumns returns the standard select columns of the model.
	GetDefaultColumns() []string

	// GetSelectedColumns returns the set of columns to be fetched.
	GetSelectedColumns() []string

	// IsDistinct returns true if SELECT DISTINCT is requested.
	IsDistinct() bool

	// GetAggregate returns the select aggregate function state.
	GetAggregate() *AggregateState

	// GetWheres returns all applied query WHERE conditions.
	GetWheres() []Where

	// GetOrders returns all applied sorting ORDER BY clauses.
	GetOrders() []Order

	// GetGroups returns the list of columns to GROUP BY.
	GetGroups() []string

	// GetHavings returns all applied grouped HAVING conditions.
	GetHavings() []Where

	// GetLimit returns the maximum number of rows to return.
	GetLimit() (uint64, bool)

	// GetOffset returns the number of rows to skip before returning.
	GetOffset() (uint64, bool)

	// GetColumnCastAndTypedSlice retrieves the SQL cast suffix, converts
	// the untyped slice of values to its strongly-typed counterpart, and
	// returns whether it is an array.
	GetColumnCastAndTypedSlice(
		colName string,
		slice []any,
	) (cast string, typedSlice any, isArray bool)
}
