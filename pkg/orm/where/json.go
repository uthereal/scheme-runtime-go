package where

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// JsonWhere represents a JSON key comparison
// (e.g., column->>'key' = value).
type JsonWhere struct {
	Column   string
	Key      string
	Operator contract.ComparisonOperator
	Value    any
	Boolean  contract.BooleanOperator
}

// GetColumn returns the name of the base JSON column.
func (w JsonWhere) GetColumn() string {
	return w.Column
}

// GetKey returns the JSON property path/key string being evaluated.
func (w JsonWhere) GetKey() string {
	return w.Key
}

// GetOperator returns the comparison operator.
func (w JsonWhere) GetOperator() contract.ComparisonOperator {
	return w.Operator
}

// GetValue returns the literal comparison value.
func (w JsonWhere) GetValue() any {
	return w.Value
}

// GetBoolean returns the conjunction operator ("AND" / "OR") linking
// this condition.
func (w JsonWhere) GetBoolean() contract.BooleanOperator {
	return w.Boolean
}

// WithBoolean returns a copy of the condition with the new conjunction
// operator.
func (w JsonWhere) WithBoolean(op contract.BooleanOperator) contract.Where {
	w.Boolean = op
	return w
}
