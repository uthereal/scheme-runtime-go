package column

import (
	"cmp"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// NumericColumn represents a strongly-typed column schema for numeric values
// (integers and floats).
type NumericColumn[Model any, Type cmp.Ordered] struct {
	Column[Model, Type]
}

// Gt returns a greater-than comparison condition (e.g., column > value).
func (c NumericColumn[Model, Type]) Gt(val Type) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpGreaterThan,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// Gte returns a greater-than-or-equal-to comparison condition
// (e.g., column >= value).
func (c NumericColumn[Model, Type]) Gte(val Type) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpGreaterThanOrEqual,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// Lt returns a less-than comparison condition (e.g., column < value).
func (c NumericColumn[Model, Type]) Lt(val Type) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpLessThan,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// Lte returns a less-than-or-equal-to comparison condition
// (e.g., column <= value).
func (c NumericColumn[Model, Type]) Lte(val Type) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpLessThanOrEqual,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}
