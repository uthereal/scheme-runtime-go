package where

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// NullWhere represents a nullability check condition (e.g., column IS NULL
// or column IS NOT NULL).
type NullWhere struct {
	Column  string
	Not     bool
	Boolean contract.BooleanOperator
}

// GetColumn returns the name of the column being evaluated.
func (w NullWhere) GetColumn() string {
	return w.Column
}

// IsNot returns true if the check is "IS NOT NULL".
func (w NullWhere) IsNot() bool {
	return w.Not
}

// GetBoolean returns the conjunction operator linking this condition.
func (w NullWhere) GetBoolean() contract.BooleanOperator {
	return w.Boolean
}

// WithBoolean returns a copy of the condition with the new conjunction
// operator.
func (w NullWhere) WithBoolean(op contract.BooleanOperator) contract.Where {
	w.Boolean = op
	return w
}
