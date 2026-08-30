package column

import (
	"encoding/json"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// ArrayColumn represents a strongly-typed column schema for database arrays.
type ArrayColumn[Model any, Elem any] struct {
	Column[Model, []Elem]
}

// ToTypedSlice converts a generic slice of []any to a strongly-typed string
// representation.
func (ac ArrayColumn[Model, Elem]) ToTypedSlice(slice []any) any {
	res := make([]string, len(slice))
	for i, v := range slice {
		if v == nil {
			res[i] = "{}"
			continue
		}
		var target []Elem
		switch val := v.(type) {
		case []Elem:
			target = val
		case *[]Elem:
			if val != nil {
				target = *val
			}
		}
		if len(target) == 0 {
			res[i] = "{}"
			continue
		}
		bytes, err := json.Marshal(target)
		if err != nil {
			res[i] = "{}"
			continue
		}
		str := string(bytes)
		if len(str) >= 2 && str[0] == '[' && str[len(str)-1] == ']' {
			str = "{" + str[1:len(str)-1] + "}"
		}
		res[i] = str
	}
	return res
}

// PostgresCast returns the PostgreSQL cast suffix for ArrayColumn.
func (ac ArrayColumn[Model, Elem]) PostgresCast() string {
	return "::text[]"
}

// IsArray returns true for ArrayColumn.
func (ac ArrayColumn[Model, Elem]) IsArray() bool {
	return true
}

// Contains returns a condition checking if the array contains all
// specified values.
func (ac ArrayColumn[Model, Elem]) Contains(vals []Elem) contract.Where {
	return where.BasicWhere{
		Column:   ac.Name,
		Operator: contract.OpArrayContains,
		Value:    vals,
		Boolean:  contract.BoolAnd,
	}
}

// Overlaps returns a condition checking if the array shares any
// elements with the specified values.
func (ac ArrayColumn[Model, Elem]) Overlaps(vals []Elem) contract.Where {
	return where.BasicWhere{
		Column:   ac.Name,
		Operator: contract.OpArrayOverlap,
		Value:    vals,
		Boolean:  contract.BoolAnd,
	}
}

// ContainedBy returns a condition checking if the array is contained by
// the specified values.
func (ac ArrayColumn[Model, Elem]) ContainedBy(vals []Elem) contract.Where {
	return where.BasicWhere{
		Column:   ac.Name,
		Operator: contract.OpArrayContainedBy,
		Value:    vals,
		Boolean:  contract.BoolAnd,
	}
}

// Concat returns a condition checking if concatenating elements matches the
// given array.
func (ac ArrayColumn[Model, Elem]) Concat(vals []Elem) contract.Where {
	return where.BasicWhere{
		Column:   ac.Name,
		Operator: contract.OpArrayConcat,
		Value:    vals,
		Boolean:  contract.BoolAnd,
	}
}
