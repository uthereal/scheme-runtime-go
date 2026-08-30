package contract

// Where represents any valid SQL WHERE condition clause.
type Where interface {
	// GetBoolean returns the conjunction operator linking this condition to
	// the parent.
	GetBoolean() BooleanOperator

	// WithBoolean returns a copy of this Where condition with the new
	// conjunction operator.
	WithBoolean(op BooleanOperator) Where
}

// BasicWhere represents a standard column-value comparison condition.
type BasicWhere interface {
	Where

	// GetColumn returns the name of the column being evaluated.
	GetColumn() string

	// GetOperator returns the comparison operator.
	GetOperator() ComparisonOperator

	// GetValue returns the literal value being compared against.
	GetValue() any

	// GetBoolean returns the conjunction operator linking this condition.
	GetBoolean() BooleanOperator
}

// BetweenWhere represents a range-boundary comparison condition.
type BetweenWhere interface {
	Where

	// GetColumn returns the name of the column being evaluated.
	GetColumn() string

	// GetMin returns the lower bound value of the range.
	GetMin() any

	// GetMax returns the upper bound value of the range.
	GetMax() any

	// IsNot returns true if the BETWEEN check is negated.
	IsNot() bool

	// GetBoolean returns the conjunction operator linking this condition.
	GetBoolean() BooleanOperator
}

// ColumnWhere represents a comparison condition between two columns.
type ColumnWhere interface {
	Where

	// GetFirst returns the name of the first column being compared.
	GetFirst() string

	// GetOperator returns the comparison operator.
	GetOperator() ComparisonOperator

	// GetSecond returns the name of the second column being compared.
	GetSecond() string

	// GetBoolean returns the conjunction operator linking this condition.
	GetBoolean() BooleanOperator
}

// DateWhere represents a date/time component extraction and comparison.
type DateWhere interface {
	Where

	// GetType returns the date component/extraction type.
	GetType() string

	// GetColumn returns the name of the datetime column.
	GetColumn() string

	// GetOperator returns the comparison operator.
	GetOperator() ComparisonOperator

	// GetValue returns the literal value being compared against.
	GetValue() any

	// GetBoolean returns the conjunction operator linking this condition.
	GetBoolean() BooleanOperator
}

// ExistsWhere represents a sub-query existence check condition.
type ExistsWhere interface {
	Where

	// GetQuery returns the sub-query provider evaluated for existence.
	GetQuery() QueryStateProvider

	// IsNot returns true if the EXISTS check is negated.
	IsNot() bool

	// GetBoolean returns the conjunction operator linking this condition.
	GetBoolean() BooleanOperator
}

// InWhere represents a set-membership evaluation condition.
type InWhere interface {
	Where

	// GetColumn returns the name of the column being evaluated.
	GetColumn() string

	// GetValues returns the slice of literal values to check.
	GetValues() []any

	// IsNot returns true if the IN check is negated.
	IsNot() bool

	// GetBoolean returns the conjunction operator linking this condition.
	GetBoolean() BooleanOperator
}

// JsonWhere represents a JSON key comparison condition.
type JsonWhere interface {
	Where

	// GetColumn returns the name of the base JSON column.
	GetColumn() string

	// GetKey returns the JSON property path/key string being evaluated.
	GetKey() string

	// GetOperator returns the comparison operator.
	GetOperator() ComparisonOperator

	// GetValue returns the literal comparison value.
	GetValue() any

	// GetBoolean returns the conjunction operator linking this condition.
	GetBoolean() BooleanOperator
}

// NestedWhere represents a grouped/nested set of query conditions.
type NestedWhere interface {
	Where

	// GetQuery returns the sub-query provider with nested conditions.
	GetQuery() QueryStateProvider

	// IsNegated returns true if the nested conditions are negated.
	IsNegated() bool

	// GetBoolean returns the conjunction operator linking this condition.
	GetBoolean() BooleanOperator
}

// NullWhere represents a nullability check condition.
type NullWhere interface {
	Where

	// GetColumn returns the name of the column being evaluated.
	GetColumn() string

	// IsNot returns true if the NULL check is negated.
	IsNot() bool

	// GetBoolean returns the conjunction operator linking this condition.
	GetBoolean() BooleanOperator
}

// RawWhere represents a raw SQL condition statement with binding arguments.
type RawWhere interface {
	Where

	// GetSql returns the raw SQL fragment of the condition.
	GetSql() string

	// GetBindings returns the bound parameters for the raw SQL query.
	GetBindings() []any

	// GetBoolean returns the conjunction operator linking this condition.
	GetBoolean() BooleanOperator
}

// SubQueryWhere represents a column comparison against a scalar sub-query.
type SubQueryWhere interface {
	Where

	// GetColumn returns the name of the column being evaluated.
	GetColumn() string

	// GetOperator returns the comparison operator.
	GetOperator() ComparisonOperator

	// GetQuery returns the sub-query provider yielding the scalar result.
	GetQuery() QueryStateProvider

	// GetBoolean returns the conjunction operator linking this condition.
	GetBoolean() BooleanOperator
}
