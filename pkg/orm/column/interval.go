package column

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/order"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// IntervalColumn represents a strongly-typed column schema for
// PostgreSQL interval (pgtype.Interval) values.
type IntervalColumn[Model any, Type pgtype.Interval] struct {
	Column[Model, Type]
}

// Gt returns a greater-than comparison condition (e.g., column > value).
func (c IntervalColumn[Model, Type]) Gt(val Type) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpGreaterThan,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// Gte returns a greater-than-or-equal-to comparison condition (e.g., column >= value).
func (c IntervalColumn[Model, Type]) Gte(val Type) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpGreaterThanOrEqual,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// Lt returns a less-than comparison condition (e.g., column < value).
func (c IntervalColumn[Model, Type]) Lt(val Type) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpLessThan,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// Lte returns a less-than-or-equal-to comparison condition (e.g., column <= value).
func (c IntervalColumn[Model, Type]) Lte(val Type) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpLessThanOrEqual,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// NullableIntervalColumn represents a strongly-typed column schema for
// nullable PostgreSQL interval (pgtype.Interval) fields.
type NullableIntervalColumn[Model any, Type pgtype.Interval] struct {
	IntervalColumn[Model, Type]
}

// IsNull returns a condition checking if this column is NULL.
func (c NullableIntervalColumn[Model, Type]) IsNull() contract.Where {
	return where.NullWhere{
		Column:  c.Name,
		Not:     false,
		Boolean: contract.BoolAnd,
	}
}

// IsNotNull returns a condition checking if this column is NOT NULL.
func (c NullableIntervalColumn[Model, Type]) IsNotNull() contract.Where {
	return where.NullWhere{
		Column:  c.Name,
		Not:     true,
		Boolean: contract.BoolAnd,
	}
}

// EqPtr returns an equality check condition against a pointer value.
func (c NullableIntervalColumn[Model, Type]) EqPtr(val *Type) contract.Where {
	return c.Eq(*val)
}

// NeqPtr returns an inequality check condition against a pointer value.
func (c NullableIntervalColumn[Model, Type]) NeqPtr(val *Type) contract.Where {
	return c.Neq(*val)
}

// GtPtr returns a greater-than comparison condition against a pointer value.
func (c NullableIntervalColumn[Model, Type]) GtPtr(val *Type) contract.Where {
	return c.Gt(*val)
}

// GtePtr returns a greater-than-or-equal-to comparison condition against a pointer value.
func (c NullableIntervalColumn[Model, Type]) GtePtr(val *Type) contract.Where {
	return c.Gte(*val)
}

// LtPtr returns a less-than comparison condition against a pointer value.
func (c NullableIntervalColumn[Model, Type]) LtPtr(val *Type) contract.Where {
	return c.Lt(*val)
}

// LtePtr returns a less-than-or-equal-to comparison condition against a pointer value.
func (c NullableIntervalColumn[Model, Type]) LtePtr(val *Type) contract.Where {
	return c.Lte(*val)
}

// AscNullsFirst returns an ascending sort order with NULL values positioned first.
func (c NullableIntervalColumn[Model, Type]) AscNullsFirst() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortAsc,
		Nulls:     contract.NullsFirst,
	}
}

// AscNullsLast returns an ascending sort order with NULL values positioned last.
func (c NullableIntervalColumn[Model, Type]) AscNullsLast() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortAsc,
		Nulls:     contract.NullsLast,
	}
}

// DescNullsFirst returns a descending sort order with NULL values positioned first.
func (c NullableIntervalColumn[Model, Type]) DescNullsFirst() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortDesc,
		Nulls:     contract.NullsFirst,
	}
}

// DescNullsLast returns a descending sort order with NULL values positioned last.
func (c NullableIntervalColumn[Model, Type]) DescNullsLast() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortDesc,
		Nulls:     contract.NullsLast,
	}
}
