package where

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// BetweenWhere represents a range-boundary comparison
// condition (e.g., column BETWEEN min AND max).
type BetweenWhere struct {
	Column  string
	Min     any
	Max     any
	Not     bool
	Boolean contract.BooleanOperator
}

// GetColumn returns the name of the column being evaluated.
func (w BetweenWhere) GetColumn() string {
	return w.Column
}

// GetMin returns the lower bound value of the range.
func (w BetweenWhere) GetMin() any {
	return w.Min
}

// GetMax returns the upper bound value of the range.
func (w BetweenWhere) GetMax() any {
	return w.Max
}

// IsNot returns true if the check is "NOT BETWEEN".
func (w BetweenWhere) IsNot() bool {
	return w.Not
}

// GetBoolean returns the conjunction operator linking this condition.
func (w BetweenWhere) GetBoolean() contract.BooleanOperator {
	return w.Boolean
}

// WithBoolean returns a copy of the condition with the new conjunction
// operator.
func (w BetweenWhere) WithBoolean(op contract.BooleanOperator) contract.Where {
	w.Boolean = op
	return w
}
