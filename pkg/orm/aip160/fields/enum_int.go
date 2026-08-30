package fields

import (
	"fmt"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	ormaip160 "github.com/uthereal/scheme-runtime-go/pkg/orm/aip160"
	luciaip160 "go.chromium.org/luci/common/data/aip160"
	"golang.org/x/exp/constraints"
)

// integerEnumField is a filter field for integer-based enums.
type integerEnumField[T constraints.Integer, R any] struct {
	subject contract.Aip160EnumSubject[T, R]
}

// NewIntegerEnumField creates a new filter field for integer-based enums.
func NewIntegerEnumField[T constraints.Integer, R any](
	subject contract.Aip160EnumSubject[T, R],
) contract.Aip160Field[R] {
	return &integerEnumField[T, R]{subject: subject}
}

// Type returns the AIP-160 field type for this field.
func (f *integerEnumField[T, R]) Type() contract.Aip160FieldType {
	return contract.TypeEnum
}

// ValidateArg validates whether the provided argument is a valid integer.
func (f *integerEnumField[T, R]) ValidateArg(arg *luciaip160.Arg) error {
	_, err := ormaip160.CoerceArgToIntConstant(arg)
	return err
}

// BuildCondition builds the query condition for the integer enum field.
func (f *integerEnumField[T, R]) BuildCondition(
	op contract.Aip160Operator,
	arg *luciaip160.Arg,
) (R, error) {
	var zero R
	val, err := ormaip160.CoerceArgToIntConstant(arg)
	if err != nil {
		return zero, err
	}
	switch op {
	case contract.OpEq, contract.OpDefault:
		return f.subject.Eq(T(val)), nil
	case contract.OpNotEq:
		return f.subject.Neq(T(val)), nil
	default:
		return zero, fmt.Errorf(
			"unsupported operator %q for integer enum field",
			op,
		)
	}
}
