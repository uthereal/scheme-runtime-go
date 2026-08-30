package fields

import (
	"fmt"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	luciaip160 "go.chromium.org/luci/common/data/aip160"
)

// relatedField represents a filter field for relation columns.
type relatedField[R any] struct {
	baseField contract.Aip160Field[R]
	mapper    func(childCond R) (R, error)
}

// NewRelatedField creates a new field that delegates operator evaluation
// to a base field and lifts the result via a mapper function.
func NewRelatedField[R any](
	baseField contract.Aip160Field[R],
	mapper func(childCond R) (R, error),
) contract.Aip160Field[R] {
	return &relatedField[R]{
		baseField: baseField,
		mapper:    mapper,
	}
}

// Type returns the field type.
func (f *relatedField[R]) Type() contract.Aip160FieldType {
	return f.baseField.Type()
}

// ValidateArg validates that the provided argument is valid for the base
// field.
func (f *relatedField[R]) ValidateArg(arg *luciaip160.Arg) error {
	return f.baseField.ValidateArg(arg)
}

// BuildCondition builds an AIP-160 filter condition by delegating to the
// base field and mapping the result.
func (f *relatedField[R]) BuildCondition(
	op contract.Aip160Operator,
	arg *luciaip160.Arg,
) (R, error) {
	var zero R
	childCond, err := f.baseField.BuildCondition(op, arg)
	if err != nil {
		return zero, fmt.Errorf(
			"failed to build child condition -> %w",
			err,
		)
	}
	return f.mapper(childCond)
}
