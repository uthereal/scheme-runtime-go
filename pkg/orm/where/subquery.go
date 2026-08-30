package where

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// SubQueryWhere represents a column comparison against a scalar sub-query
// result (e.g., column = (SELECT ...)).
type SubQueryWhere struct {
	Column   string
	Operator contract.ComparisonOperator
	Query    contract.QueryStateProvider
	Boolean  contract.BooleanOperator
}

// GetColumn returns the name of the column being evaluated.
func (w SubQueryWhere) GetColumn() string {
	return w.Column
}

// GetOperator returns the comparison operator.
func (w SubQueryWhere) GetOperator() contract.ComparisonOperator {
	return w.Operator
}

// GetQuery returns the sub-query provider yielding the scalar result.
func (w SubQueryWhere) GetQuery() contract.QueryStateProvider {
	return w.Query
}

// GetBoolean returns the conjunction operator linking this condition.
func (w SubQueryWhere) GetBoolean() contract.BooleanOperator {
	return w.Boolean
}

// WithBoolean returns a copy of the condition with the new conjunction
// operator.
func (w SubQueryWhere) WithBoolean(op contract.BooleanOperator) contract.Where {
	w.Boolean = op
	return w
}
