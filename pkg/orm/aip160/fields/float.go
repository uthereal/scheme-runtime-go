package fields

import (
	"fmt"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	ormaip160 "github.com/uthereal/scheme-runtime-go/pkg/orm/aip160"
	luciaip160 "go.chromium.org/luci/common/data/aip160"
	"golang.org/x/exp/constraints"
)

// floatField is a filter field for float-like columns.
type floatField[T constraints.Float, R any] struct {
	subject contract.Aip160NumericSubject[T, R]
}

// NewFloatField creates a new filter field for float-like columns.
func NewFloatField[T constraints.Float, R any](
	subject contract.Aip160NumericSubject[T, R],
) contract.Aip160Field[R] {
	return &floatField[T, R]{subject: subject}
}

// Type returns the AIP-160 field type for this field.
func (f *floatField[T, R]) Type() contract.Aip160FieldType {
	return contract.TypeFloat
}

// ValidateArg validates whether the provided argument is a valid float.
func (f *floatField[T, R]) ValidateArg(arg *luciaip160.Arg) error {
	_, err := ormaip160.CoerceArgToFloatConstant(arg)
	return err
}

// BuildCondition builds the query condition for the float field.
func (f *floatField[T, R]) BuildCondition(
	op contract.Aip160Operator,
	arg *luciaip160.Arg,
) (R, error) {
	var zero R
	val, err := ormaip160.CoerceArgToFloatConstant(arg)
	if err != nil {
		return zero, err
	}
	switch op {
	case contract.OpEq, contract.OpDefault:
		return f.subject.Eq(T(val)), nil
	case contract.OpNotEq:
		return f.subject.Neq(T(val)), nil
	case contract.OpGt:
		return f.subject.Gt(T(val)), nil
	case contract.OpGtEq:
		return f.subject.Gte(T(val)), nil
	case contract.OpLt:
		return f.subject.Lt(T(val)), nil
	case contract.OpLtEq:
		return f.subject.Lte(T(val)), nil
	default:
		return zero, fmt.Errorf(
			"unsupported operator %q for float field",
			op,
		)
	}
}
