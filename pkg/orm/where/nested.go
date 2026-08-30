package where

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// NestedWhere represents a grouped/nested set of query conditions
// enclosed in parentheses.
type NestedWhere struct {
	Query   contract.QueryStateProvider
	Not     bool
	Boolean contract.BooleanOperator
}

// GetQuery returns the sub-query provider containing the nested
// conditions.
func (w NestedWhere) GetQuery() contract.QueryStateProvider {
	return w.Query
}

// IsNegated returns true if the check is negated.
func (w NestedWhere) IsNegated() bool {
	return w.Not
}

// GetBoolean returns the conjunction operator linking this nested
// group to the parent query.
func (w NestedWhere) GetBoolean() contract.BooleanOperator {
	return w.Boolean
}

// WithBoolean returns a copy of the condition with the new conjunction
// operator.
func (w NestedWhere) WithBoolean(op contract.BooleanOperator) contract.Where {
	w.Boolean = op
	return w
}
