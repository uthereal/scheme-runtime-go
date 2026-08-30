package where

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// BasicWhere represents a standard column-value comparison
// condition (e.g., column = value).
type BasicWhere struct {
	Column   string
	Operator contract.ComparisonOperator
	Value    any
	Boolean  contract.BooleanOperator
}

// GetColumn returns the name of the column being evaluated.
func (w BasicWhere) GetColumn() string {
	return w.Column
}

// GetOperator returns the comparison operator.
func (w BasicWhere) GetOperator() contract.ComparisonOperator {
	return w.Operator
}

// GetValue returns the literal value being compared against.
func (w BasicWhere) GetValue() any {
	return w.Value
}

// GetBoolean returns the conjunction operator linking this condition.
func (w BasicWhere) GetBoolean() contract.BooleanOperator {
	return w.Boolean
}

// WithBoolean returns a copy of the condition with the new conjunction
// operator.
func (w BasicWhere) WithBoolean(op contract.BooleanOperator) contract.Where {
	w.Boolean = op
	return w
}
