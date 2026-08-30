package where

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// ColumnWhere represents a comparison condition between two
// columns (e.g., columnA = columnB).
type ColumnWhere struct {
	First    string
	Operator contract.ComparisonOperator
	Second   string
	Boolean  contract.BooleanOperator
}

// GetFirst returns the name of the first column being compared.
func (w ColumnWhere) GetFirst() string {
	return w.First
}

// GetOperator returns the comparison operator between the two columns.
func (w ColumnWhere) GetOperator() contract.ComparisonOperator {
	return w.Operator
}

// GetSecond returns the name of the second column being compared.
func (w ColumnWhere) GetSecond() string {
	return w.Second
}

// GetBoolean returns the conjunction operator linking this condition.
func (w ColumnWhere) GetBoolean() contract.BooleanOperator {
	return w.Boolean
}

// WithBoolean returns a copy of the condition with the new conjunction
// operator.
func (w ColumnWhere) WithBoolean(op contract.BooleanOperator) contract.Where {
	w.Boolean = op
	return w
}
