package column

import (
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/order"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// ByteColumn represents a strongly-typed column schema for PostgreSQL
// bytea binary fields.
type ByteColumn[Model any] struct {
	Column[Model, []byte]
}

// Compile-time interface assertion ensuring ByteColumn implements contract.Column.
var _ contract.Column[any] = (*ByteColumn[any])(nil)

// PostgresCast returns the PostgreSQL cast suffix for bytea array.
func (c ByteColumn[Model]) PostgresCast() string {
	return "::bytea[]"
}

// NullableByteColumn represents a strongly-typed column schema for nullable
// PostgreSQL bytea binary fields.
type NullableByteColumn[Model any] struct {
	ByteColumn[Model]
}

// Compile-time interface assertion ensuring NullableByteColumn implements contract.Column.
var _ contract.Column[any] = (*NullableByteColumn[any])(nil)

// IsNull returns a condition checking if this column is NULL.
func (c NullableByteColumn[Model]) IsNull() contract.Where {
	return where.NullWhere{
		Column:  c.Name,
		Not:     false,
		Boolean: contract.BoolAnd,
	}
}

// IsNotNull returns a condition checking if this column is NOT NULL.
func (c NullableByteColumn[Model]) IsNotNull() contract.Where {
	return where.NullWhere{
		Column:  c.Name,
		Not:     true,
		Boolean: contract.BoolAnd,
	}
}

// EqPtr returns an equality check condition against a pointer value.
func (c NullableByteColumn[Model]) EqPtr(val *[]byte) contract.Where {
	return c.Eq(*val)
}

// NeqPtr returns an inequality check condition against a pointer value.
func (c NullableByteColumn[Model]) NeqPtr(val *[]byte) contract.Where {
	return c.Neq(*val)
}

// AscNullsFirst returns an ascending sort order with NULL values positioned first.
func (c NullableByteColumn[Model]) AscNullsFirst() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortAsc,
		Nulls:     contract.NullsFirst,
	}
}

// AscNullsLast returns an ascending sort order with NULL values positioned last.
func (c NullableByteColumn[Model]) AscNullsLast() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortAsc,
		Nulls:     contract.NullsLast,
	}
}

// DescNullsFirst returns a descending sort order with NULL values positioned first.
func (c NullableByteColumn[Model]) DescNullsFirst() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortDesc,
		Nulls:     contract.NullsFirst,
	}
}

// DescNullsLast returns a descending sort order with NULL values positioned last.
func (c NullableByteColumn[Model]) DescNullsLast() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortDesc,
		Nulls:     contract.NullsLast,
	}
}
