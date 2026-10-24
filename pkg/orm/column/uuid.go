package column

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

// UUIDColumn represents a strongly-typed column schema for UUID fields.
type UUIDColumn[Model any] struct {
	Column[Model, uuid.UUID]
}

// PostgresCast returns the PostgreSQL cast suffix for UUID array.
func (c UUIDColumn[Model]) PostgresCast() string {
	return "::uuid[]"
}

// ToTypedSlice converts an untyped slice of values to a strongly-typed slice.
func (c UUIDColumn[Model]) ToTypedSlice(slice []any) any {
	if len(slice) == 0 {
		return []uuid.UUID{}
	}

	isPointer := false
	for _, v := range slice {
		if v != nil {
			isPointer = reflect.TypeOf(v).Kind() == reflect.Pointer
			break
		}
	}

	if isPointer {
		res := make([]*uuid.UUID, len(slice))
		for i, v := range slice {
			if v == nil {
				continue
			}
			switch val := v.(type) {
			case *uuid.UUID:
				res[i] = val
			case uuid.UUID:
				res[i] = &val
			case string:
				if parsed, err := uuid.Parse(val); err == nil {
					res[i] = &parsed
				}
			case *string:
				if val != nil {
					if parsed, err := uuid.Parse(*val); err == nil {
						res[i] = &parsed
					}
				}
			case fmt.Stringer:
				if parsed, err := uuid.Parse(val.String()); err == nil {
					res[i] = &parsed
				}
			default:
				rv := reflect.ValueOf(v)
				if rv.Kind() == reflect.Pointer {
					if rv.IsNil() {
						continue
					}
					elem := rv.Elem().Interface()
					if stringer, ok := elem.(fmt.Stringer); ok {
						if parsed, err := uuid.Parse(
							stringer.String(),
						); err == nil {
							res[i] = &parsed
						}
						continue
					}
				}
				if parsed, err := uuid.Parse(fmt.Sprint(v)); err == nil {
					res[i] = &parsed
				}
			}
		}
		return res
	}

	res := make([]uuid.UUID, len(slice))
	for i, v := range slice {
		if v == nil {
			continue
		}
		switch val := v.(type) {
		case uuid.UUID:
			res[i] = val
		case *uuid.UUID:
			if val != nil {
				res[i] = *val
			}
		case string:
			if parsed, err := uuid.Parse(val); err == nil {
				res[i] = parsed
			}
		case *string:
			if val != nil {
				if parsed, err := uuid.Parse(*val); err == nil {
					res[i] = parsed
				}
			}
		case fmt.Stringer:
			if parsed, err := uuid.Parse(val.String()); err == nil {
				res[i] = parsed
			}
		default:
			if parsed, err := uuid.Parse(fmt.Sprint(v)); err == nil {
				res[i] = parsed
			}
		}
	}
	return res
}
