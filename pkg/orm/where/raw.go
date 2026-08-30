package where

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// RawWhere represents a raw SQL condition statement with its associated
// binding arguments.
type RawWhere struct {
	Sql      string
	Bindings []any
	Boolean  contract.BooleanOperator
}

// GetSql returns the raw SQL fragment of the condition.
func (w RawWhere) GetSql() string {
	return w.Sql
}

// GetBindings returns the bound parameters for the raw SQL query string
// placeholders.
func (w RawWhere) GetBindings() []any {
	return w.Bindings
}

// GetBoolean returns the conjunction operator linking this condition.
func (w RawWhere) GetBoolean() contract.BooleanOperator {
	return w.Boolean
}

// WithBoolean returns a copy of the condition with the new conjunction
// operator.
func (w RawWhere) WithBoolean(op contract.BooleanOperator) contract.Where {
	w.Boolean = op
	return w
}
