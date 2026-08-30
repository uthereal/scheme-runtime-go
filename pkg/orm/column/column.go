// Package column defines the generic, strongly-typed column schemas
// used by developers to build type-safe query constraints.
package column

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/lo"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/order"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// Column represents a base database column schema mapped to a specific
// Model and field type.
type Column[Model any, Type any] struct {
	Name string
}

// ColumnName returns the physical name of the column in the database table.
func (c Column[Model, Type]) ColumnName() string {
	return c.Name
}

// PostgresCast returns the correct Postgres cast suffix based on Type.
func (c Column[Model, Type]) PostgresCast() string {
	var zero Type
	switch any(zero).(type) {
	case string:
		return "::text[]"
	case int, int32, *int, *int32:
		return "::integer[]"
	case int64, *int64:
		return "::bigint[]"
	case bool, *bool:
		return "::boolean[]"
	case float64, *float64:
		return "::double precision[]"
	case time.Time, *time.Time:
		return "::timestamp with time zone[]"
	case time.Duration, *time.Duration:
		return "::interval[]"
	case pgtype.Numeric, *pgtype.Numeric:
		return "::numeric[]"
	case pgtype.Point, *pgtype.Point:
		return "::point[]"
	default:
		return "::text[]"
	}
}

// ToTypedSlice converts a generic []any slice to a strongly-typed
// []Type or []*Type slice.
func (c Column[Model, Type]) ToTypedSlice(slice []any) any {
	var isPointer bool
	for _, v := range slice {
		if v == nil {
			continue
		}
		_, isPointer = v.(*Type)
		break
	}

	if isPointer {
		res := make([]*Type, len(slice))
		for i, v := range slice {
			if v == nil {
				continue
			}
			res[i] = v.(*Type)
		}
		return res
	}

	res := make([]Type, len(slice))
	for i, v := range slice {
		if v == nil {
			continue
		}
		switch val := v.(type) {
		case Type:
			res[i] = val
		case *Type:
			if val == nil {
				continue
			}
			res[i] = *val
		default:
			res[i] = v.(Type)
		}
	}
	return res
}

// IsArray returns false by default for standard scalar columns.
func (c Column[Model, Type]) IsArray() bool {
	return false
}

// BindType returns a zero value of Type to satisfy type-binding
// contract interfaces.
func (c Column[Model, Type]) BindType() Type {
	var zero Type
	return zero
}

// Asc returns an ascending order clause for this column.
func (c Column[Model, Type]) Asc() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortAsc,
	}
}

// Desc returns a descending order clause for this column.
func (c Column[Model, Type]) Desc() contract.Order {
	return order.ColumnOrder{
		Column:    c.Name,
		Direction: contract.SortDesc,
	}
}

// Eq returns a standard equality check condition (e.g., column = value).
func (c Column[Model, Type]) Eq(val Type) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpEqual,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// Neq returns a standard inequality check condition (e.g., column != value).
func (c Column[Model, Type]) Neq(val Type) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpNotEqual,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// In returns a set-membership evaluation condition
// (e.g., column IN (values)).
func (c Column[Model, Type]) In(vals ...Type) contract.Where {
	return where.InWhere{
		Column:  c.Name,
		Values:  lo.ToAnySlice(vals),
		Not:     false,
		Boolean: contract.BoolAnd,
	}
}

// NotIn returns a negative set-membership evaluation condition
// (e.g., column NOT IN (values)).
func (c Column[Model, Type]) NotIn(vals ...Type) contract.Where {
	return where.InWhere{
		Column:  c.Name,
		Values:  lo.ToAnySlice(vals),
		Not:     true,
		Boolean: contract.BoolAnd,
	}
}

// InQuery returns a set-membership evaluation condition against a
// sub-query.
func (c Column[Model, Type]) InQuery(subQuery any) contract.Where {
	return where.SubQueryWhere{
		Column:   c.Name,
		Operator: contract.OpIn,
		Query:    subQuery.(contract.QueryStateProvider),
		Boolean:  contract.BoolAnd,
	}
}

// NotInQuery returns a negative set-membership evaluation condition
// against a sub-query.
func (c Column[Model, Type]) NotInQuery(subQuery any) contract.Where {
	return where.SubQueryWhere{
		Column:   c.Name,
		Operator: contract.OpNotIn,
		Query:    subQuery.(contract.QueryStateProvider),
		Boolean:  contract.BoolAnd,
	}
}

// Between returns a range comparison condition
// (e.g., column BETWEEN min AND max).
func (c Column[Model, Type]) Between(min Type, max Type) contract.Where {
	return where.BetweenWhere{
		Column:  c.Name,
		Min:     min,
		Max:     max,
		Not:     false,
		Boolean: contract.BoolAnd,
	}
}

// NotBetween returns a negative range comparison condition
// (e.g., column NOT BETWEEN min AND max).
func (c Column[Model, Type]) NotBetween(min Type, max Type) contract.Where {
	return where.BetweenWhere{
		Column:  c.Name,
		Min:     min,
		Max:     max,
		Not:     true,
		Boolean: contract.BoolAnd,
	}
}
