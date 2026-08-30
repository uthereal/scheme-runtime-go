package column

import (
	"fmt"
	"reflect"
)

// UUIDColumn represents a strongly-typed column schema for UUID fields.
type UUIDColumn[Model any] struct {
	Column[Model, string]
}

// PostgresCast returns the PostgreSQL cast suffix for UUID array.
func (c UUIDColumn[Model]) PostgresCast() string {
	return "::uuid[]"
}

// ToTypedSlice converts an untyped slice of values to a strongly-typed slice.
func (c UUIDColumn[Model]) ToTypedSlice(slice []any) any {
	if len(slice) == 0 {
		return []string{}
	}

	isPointer := false
	for _, v := range slice {
		if v != nil {
			isPointer = reflect.TypeOf(v).Kind() == reflect.Pointer
			break
		}
	}

	if isPointer {
		res := make([]*string, len(slice))
		for i, v := range slice {
			if v == nil {
				continue
			}
			switch val := v.(type) {
			case *string:
				res[i] = val
			case string:
				res[i] = &val
			case fmt.Stringer:
				s := val.String()
				res[i] = &s
			default:
				rv := reflect.ValueOf(v)
				if rv.Kind() == reflect.Pointer {
					if rv.IsNil() {
						continue
					}
					elem := rv.Elem().Interface()
					if stringer, ok := elem.(fmt.Stringer); ok {
						s := stringer.String()
						res[i] = &s
						continue
					}
				}
				s := fmt.Sprint(v)
				res[i] = &s
			}
		}
		return res
	}

	res := make([]string, len(slice))
	for i, v := range slice {
		if v == nil {
			continue
		}
		switch val := v.(type) {
		case string:
			res[i] = val
		case *string:
			if val != nil {
				res[i] = *val
			}
		case fmt.Stringer:
			res[i] = val.String()
		default:
			res[i] = fmt.Sprint(v)
		}
	}
	return res
}
