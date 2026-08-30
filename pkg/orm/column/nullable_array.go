package column

import (
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/order"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// NullableArrayColumn represents a strongly-typed column schema for
// nullable array fields.
type NullableArrayColumn[Model any, Elem any] struct {
	ArrayColumn[Model, Elem]
}

// IsNull returns a condition checking if this column is NULL.
func (c NullableArrayColumn[Model, Elem]) IsNull() contract.Where {
	return where.NullWhere{
		Column:  c.Name,
		Not:     false,
		Boolean: contract.BoolAnd,
	}
}

// IsNotNull returns a condition checking if this column is NOT NULL.
func (c NullableArrayColumn[Model, Elem]) IsNotNull() contract.Where {
	return where.NullWhere{
		Column:  c.Name,
		Not:     true,
		Boolean: contract.BoolAnd,
	}
}

// EqPtr returns an equality check condition against a pointer value.
func (c NullableArrayColumn[Model, Elem]) EqPtr(
	val *[]Elem,
) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpEqual,
		Value:    *val,
		Boolean:  contract.BoolAnd,
	}
}

// NeqPtr returns an inequality check condition against a pointer value.
func (c NullableArrayColumn[Model, Elem]) NeqPtr(
	val *[]Elem,
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
func (c NullableArrayColumn[Model, Elem]) AscNullsFirst() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortAsc,
		Nulls:     contract.NullsFirst,
	}
}

// AscNullsLast returns an ascending sort order with NULL values
// positioned last.
func (c NullableArrayColumn[Model, Elem]) AscNullsLast() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortAsc,
		Nulls:     contract.NullsLast,
	}
}

// DescNullsFirst returns a descending sort order with NULL values
// positioned first.
func (c NullableArrayColumn[Model, Elem]) DescNullsFirst() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortDesc,
		Nulls:     contract.NullsFirst,
	}
}

// DescNullsLast returns a descending sort order with NULL values
// positioned last.
func (c NullableArrayColumn[Model, Elem]) DescNullsLast() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortDesc,
		Nulls:     contract.NullsLast,
	}
}
