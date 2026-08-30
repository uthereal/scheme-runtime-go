package column

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// DecimalColumn represents a strongly-typed column schema for exact
// decimal (pgtype.Numeric) values.
type DecimalColumn[Model any, Type pgtype.Numeric] struct {
	Column[Model, Type]
}

// Gt returns a greater-than comparison condition (e.g., column > value).
func (c DecimalColumn[Model, Type]) Gt(val Type) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpGreaterThan,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// Gte returns a greater-than-or-equal-to comparison condition
// (e.g., column >= value).
func (c DecimalColumn[Model, Type]) Gte(val Type) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpGreaterThanOrEqual,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// Lt returns a less-than comparison condition (e.g., column < value).
func (c DecimalColumn[Model, Type]) Lt(val Type) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpLessThan,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// Lte returns a less-than-or-equal-to comparison condition
// (e.g., column <= value).
func (c DecimalColumn[Model, Type]) Lte(val Type) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpLessThanOrEqual,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}
