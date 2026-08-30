package where

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// ExistsWhere represents a sub-query existence check
// condition (e.g., EXISTS (SELECT 1 ...)).
type ExistsWhere struct {
	Query   contract.QueryStateProvider
	Not     bool
	Boolean contract.BooleanOperator
}

// GetQuery returns the sub-query provider being evaluated for existence.
func (w ExistsWhere) GetQuery() contract.QueryStateProvider {
	return w.Query
}

// IsNot returns true if the check is "NOT EXISTS".
func (w ExistsWhere) IsNot() bool {
	return w.Not
}

// GetBoolean returns the conjunction operator linking this condition.
func (w ExistsWhere) GetBoolean() contract.BooleanOperator {
	return w.Boolean
}

// WithBoolean returns a copy of the condition with the new conjunction
// operator.
func (w ExistsWhere) WithBoolean(op contract.BooleanOperator) contract.Where {
	w.Boolean = op
	return w
}
