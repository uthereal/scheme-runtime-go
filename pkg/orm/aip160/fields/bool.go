package fields

import (
	"fmt"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	ormaip160 "github.com/uthereal/scheme-runtime-go/pkg/orm/aip160"
	luciaip160 "go.chromium.org/luci/common/data/aip160"
)

// boolField is a filter field for boolean columns.
type boolField[R any] struct {
	subject contract.Aip160BoolSubject[R]
}

// NewBoolField creates a new filter field for boolean columns.
func NewBoolField[R any](
	subject contract.Aip160BoolSubject[R],
) contract.Aip160Field[R] {
	return &boolField[R]{subject: subject}
}

// Type returns the AIP-160 field type for this field.
func (f *boolField[R]) Type() contract.Aip160FieldType {
	return contract.TypeBool
}

// ValidateArg validates whether the provided argument is a valid boolean.
func (f *boolField[R]) ValidateArg(arg *luciaip160.Arg) error {
	_, err := ormaip160.CoerceArgToBoolConstant(arg)
	return err
}

// BuildCondition builds the query condition for the boolean field.
func (f *boolField[R]) BuildCondition(
	op contract.Aip160Operator,
	arg *luciaip160.Arg,
) (R, error) {
	var zero R
	val, err := ormaip160.CoerceArgToBoolConstant(arg)
	if err != nil {
		return zero, err
	}
	switch op {
	case contract.OpEq, contract.OpDefault:
		return f.subject.Eq(val), nil
	case contract.OpNotEq:
		return f.subject.Neq(val), nil
	default:
		return zero, fmt.Errorf(
			"unsupported operator %q for boolean field",
			op,
		)
	}
}
