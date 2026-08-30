package where

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// InWhere represents a set-membership evaluation
// condition (e.g., column IN (values)).
type InWhere struct {
	Column  string
	Values  []any
	Not     bool
	Boolean contract.BooleanOperator
}

// GetColumn returns the name of the column being evaluated.
func (w InWhere) GetColumn() string {
	return w.Column
}

// GetValues returns the slice of literal values being checked
// for set membership.
func (w InWhere) GetValues() []any {
	return w.Values
}

// IsNot returns true if the check is "NOT IN".
func (w InWhere) IsNot() bool {
	return w.Not
}

// GetBoolean returns the conjunction operator linking this condition.
func (w InWhere) GetBoolean() contract.BooleanOperator {
	return w.Boolean
}

// WithBoolean returns a copy of the condition with the new conjunction
// operator.
func (w InWhere) WithBoolean(op contract.BooleanOperator) contract.Where {
	w.Boolean = op
	return w
}
