package where

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// DateWhere represents a date/time component extraction and
// comparison condition (e.g., DATE(column) = value).
type DateWhere struct {
	Type     string
	Column   string
	Operator contract.ComparisonOperator
	Value    any
	Boolean  contract.BooleanOperator
}

// GetType returns the component/type of extraction.
func (w DateWhere) GetType() string {
	return w.Type
}

// GetColumn returns the name of the column containing the datetime value.
func (w DateWhere) GetColumn() string {
	return w.Column
}

// GetOperator returns the comparison operator.
func (w DateWhere) GetOperator() contract.ComparisonOperator {
	return w.Operator
}

// GetValue returns the literal value being compared against.
func (w DateWhere) GetValue() any {
	return w.Value
}

// GetBoolean returns the conjunction operator linking this condition.
func (w DateWhere) GetBoolean() contract.BooleanOperator {
	return w.Boolean
}

// WithBoolean returns a copy of the condition with the new conjunction
// operator.
func (w DateWhere) WithBoolean(op contract.BooleanOperator) contract.Where {
	w.Boolean = op
	return w
}
