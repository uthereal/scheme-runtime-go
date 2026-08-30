package fields

import (
	"fmt"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	ormaip160 "github.com/uthereal/scheme-runtime-go/pkg/orm/aip160"
	luciaip160 "go.chromium.org/luci/common/data/aip160"
	"golang.org/x/exp/constraints"
)

// integerField represents a filter field for integer columns.
type integerField[T constraints.Integer, R any] struct {
	subject contract.Aip160NumericSubject[T, R]
}

// NewIntegerField creates a new filter field for integer-like columns.
func NewIntegerField[
	T constraints.Integer,
	R any,
](
	subject contract.Aip160NumericSubject[T, R],
) contract.Aip160Field[R] {
	return &integerField[T, R]{subject: subject}
}

// Type returns the field type.
func (f *integerField[T, R]) Type() contract.Aip160FieldType {
	return contract.TypeInt
}

// ValidateArg validates that the provided argument can be coerced to an
// integer.
func (f *integerField[T, R]) ValidateArg(arg *luciaip160.Arg) error {
	_, err := ormaip160.CoerceArgToIntConstant(arg)
	return err
}

// BuildCondition builds an AIP-160 filter condition for an integer field.
func (f *integerField[T, R]) BuildCondition(
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
			"unsupported operator %q for integer field",
			op,
		)
	}
}
