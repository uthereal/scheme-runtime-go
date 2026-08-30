package fields

import (
	"fmt"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	ormaip160 "github.com/uthereal/scheme-runtime-go/pkg/orm/aip160"
	luciaip160 "go.chromium.org/luci/common/data/aip160"
)

// stringField represents a filter field for string columns.
type stringField[T ~string, R any] struct {
	subject contract.Aip160StringSubject[T, R]
}

// NewStringField creates a new filter field for string-like columns.
func NewStringField[
	T ~string,
	R any,
](
	subject contract.Aip160StringSubject[T, R],
) contract.Aip160Field[R] {
	return &stringField[T, R]{subject: subject}
}

// Type returns the field type.
func (f *stringField[T, R]) Type() contract.Aip160FieldType {
	return contract.TypeString
}

// ValidateArg validates that the provided argument can be coerced to a
// string.
func (f *stringField[T, R]) ValidateArg(arg *luciaip160.Arg) error {
	_, err := ormaip160.CoerceArgToStringConstant(arg)
	return err
}

// BuildCondition builds an AIP-160 filter condition for a string field.
func (f *stringField[T, R]) BuildCondition(
	op contract.Aip160Operator,
	arg *luciaip160.Arg,
) (R, error) {
	var zero R
	val, err := ormaip160.CoerceArgToStringConstant(arg)
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
			"unsupported operator %q for string field",
			op,
		)
	}
}
