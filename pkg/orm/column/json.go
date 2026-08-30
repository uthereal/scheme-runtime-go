package column

import (
	"encoding/json"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// JSONColumn represents a strongly-typed column schema for JSON or
// JSONB fields.
type JSONColumn[Model any, Type any] struct {
	Column[Model, Type]
}

// ToTypedSlice converts a generic []any slice to a strongly-typed string
// slice representing JSON array elements.
func (c JSONColumn[Model, Type]) ToTypedSlice(slice []any) any {
	res := make([]string, len(slice))
	for i, v := range slice {
		if v == nil {
			res[i] = "null"
			continue
		}
		bytes, _ := json.Marshal(v)
		res[i] = string(bytes)
	}
	return res
}

// PostgresCast returns the PostgreSQL cast suffix for JSONB array.
func (c JSONColumn[Model, Type]) PostgresCast() string {
	return "::jsonb[]"
}

// IsArray returns false by default for JSONColumn.
func (c JSONColumn[Model, Type]) IsArray() bool {
	return false
}

// Contains returns a condition checking if the JSON document contains
// the specified payload.
func (c JSONColumn[Model, Type]) Contains(val any) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpJsonContains,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}

// HasKey returns a condition checking if the JSON document contains
// the specified top-level key.
func (c JSONColumn[Model, Type]) HasKey(key string) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpJsonKeyExist,
		Value:    key,
		Boolean:  contract.BoolAnd,
	}
}

// KeyEq returns a condition checking if a specific nested JSON key
// equals the given value.
func (c JSONColumn[Model, Type]) KeyEq(key string, val any) contract.Where {
	return where.JsonWhere{
		Column:   c.Name,
		Key:      key,
		Operator: contract.OpEqual,
		Value:    val,
		Boolean:  contract.BoolAnd,
	}
}
