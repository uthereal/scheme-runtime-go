package column

import (
	"cmp"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/order"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// NullableNumericColumn represents a strongly-typed column schema
// for nullable numeric fields.
type NullableNumericColumn[Model any, Type cmp.Ordered] struct {
	NumericColumn[Model, Type]
}

// IsNull returns a condition checking if this column is NULL.
func (c NullableNumericColumn[Model, Type]) IsNull() contract.Where {
	return where.NullWhere{
		Column:  c.Name,
		Not:     false,
		Boolean: contract.BoolAnd,
	}
}

// IsNotNull returns a condition checking if this column is NOT NULL.
func (c NullableNumericColumn[Model, Type]) IsNotNull() contract.Where {
	return where.NullWhere{
		Column:  c.Name,
		Not:     true,
		Boolean: contract.BoolAnd,
	}
}

// EqPtr returns an equality check condition against a pointer value.
func (c NullableNumericColumn[Model, Type]) EqPtr(
	val *Type,
) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpEqual,
		Value:    *val,
		Boolean:  contract.BoolAnd,
	}
}

// NeqPtr returns an inequality check condition against a pointer value.
func (c NullableNumericColumn[Model, Type]) NeqPtr(
	val *Type,
) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpNotEqual,
		Value:    *val,
		Boolean:  contract.BoolAnd,
	}
}

// AscNullsFirst returns an ascending sort order with NULL values
// positioned first.
func (c NullableNumericColumn[Model, Type]) AscNullsFirst() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortAsc,
		Nulls:     contract.NullsFirst,
	}
}

// AscNullsLast returns an ascending sort order with NULL values
// positioned last.
func (c NullableNumericColumn[Model, Type]) AscNullsLast() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortAsc,
		Nulls:     contract.NullsLast,
	}
}

// DescNullsFirst returns a descending sort order with NULL values
// positioned first.
func (c NullableNumericColumn[Model, Type]) DescNullsFirst() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortDesc,
		Nulls:     contract.NullsFirst,
	}
}

// DescNullsLast returns a descending sort order with NULL values
// positioned last.
func (c NullableNumericColumn[Model, Type]) DescNullsLast() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortDesc,
		Nulls:     contract.NullsLast,
	}
}

// BindType returns a zero value of Type to satisfy type-binding of the
// aggregate interface.
func (c NullableNumericColumn[Model, Type]) BindType() Type {
	var zero Type
	return zero
}
