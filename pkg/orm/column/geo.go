package column

import (
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/order"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// GeoColumn represents a strongly-typed column schema for database
// geometric types.
type GeoColumn[Model any, GeoType any] struct {
	Column[Model, GeoType]
}

// NullableGeoColumn represents a strongly-typed column schema for nullable
// geometric fields.
type NullableGeoColumn[Model any, GeoType any] struct {
	GeoColumn[Model, GeoType]
}

// Overlaps returns a condition checking if two geometric objects overlap.
func (c GeoColumn[Model, GeoType]) Overlaps(val GeoType) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpGeoOverlaps,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// Contains returns a condition checking if this geometric object contains
// another.
func (c GeoColumn[Model, GeoType]) Contains(val GeoType) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpGeoContains,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// ContainedBy returns a condition checking if this geometric object is
// contained within another.
func (c GeoColumn[Model, GeoType]) ContainedBy(val GeoType) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpGeoContainedBy,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// StrictLeft returns a condition checking if this geometric object is strictly
// to the left of another.
func (c GeoColumn[Model, GeoType]) StrictLeft(val GeoType) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpGeoStrictLeft,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// StrictRight returns a condition checking if this geometric object is
// strictly to the right of another.
func (c GeoColumn[Model, GeoType]) StrictRight(val GeoType) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpGeoStrictRight,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// Below returns a condition checking if this geometric object is below another.
func (c GeoColumn[Model, GeoType]) Below(val GeoType) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpGeoBelow,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// Above returns a condition checking if this geometric object is above another.
func (c GeoColumn[Model, GeoType]) Above(val GeoType) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpGeoAbove,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// Distance returns a condition checking the distance between two geometric
// objects.
func (c GeoColumn[Model, GeoType]) Distance(val GeoType) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpGeoDistance,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// ClosestProx returns a condition checking proximity between two geometric
// objects.
func (c GeoColumn[Model, GeoType]) ClosestProx(val GeoType) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpGeoClosestProx,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// IsNull returns a condition checking if this column is NULL.
func (c NullableGeoColumn[Model, GeoType]) IsNull() contract.Where {
	return where.NullWhere{
		Column:  c.Name,
		Not:     false,
		Boolean: contract.BoolAnd,
	}
}

// IsNotNull returns a condition checking if this column is NOT NULL.
func (c NullableGeoColumn[Model, GeoType]) IsNotNull() contract.Where {
	return where.NullWhere{
		Column:  c.Name,
		Not:     true,
		Boolean: contract.BoolAnd,
	}
}

// EqPtr returns an equality check condition against a pointer value.
func (c NullableGeoColumn[Model, GeoType]) EqPtr(
	val *GeoType,
) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpEqual,
		Value:    *val,
		Boolean:  contract.BoolAnd,
	}
}

// NeqPtr returns an inequality check condition against a pointer value.
func (c NullableGeoColumn[Model, GeoType]) NeqPtr(
	val *GeoType,
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
func (c NullableGeoColumn[Model, GeoType]) AscNullsFirst() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortAsc,
		Nulls:     contract.NullsFirst,
	}
}

// AscNullsLast returns an ascending sort order with NULL values
// positioned last.
func (c NullableGeoColumn[Model, GeoType]) AscNullsLast() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortAsc,
		Nulls:     contract.NullsLast,
	}
}

// DescNullsFirst returns a descending sort order with NULL values
// positioned first.
func (c NullableGeoColumn[Model, GeoType]) DescNullsFirst() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortDesc,
		Nulls:     contract.NullsFirst,
	}
}

// DescNullsLast returns a descending sort order with NULL values
// positioned last.
func (c NullableGeoColumn[Model, GeoType]) DescNullsLast() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortDesc,
		Nulls:     contract.NullsLast,
	}
}
