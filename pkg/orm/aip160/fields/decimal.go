package fields

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	ormaip160 "github.com/uthereal/scheme-runtime-go/pkg/orm/aip160"
	luciaip160 "go.chromium.org/luci/common/data/aip160"
)

// decimalField is a filter field for exact decimal-like (pgtype.Numeric)
// columns.
type decimalField[T pgtype.Numeric, R any] struct {
	subject contract.Aip160DecimalSubject[T, R]
}

// NewDecimalField creates a new filter field for exact decimal-like columns.
func NewDecimalField[T pgtype.Numeric, R any](
	subject contract.Aip160DecimalSubject[T, R],
) contract.Aip160Field[R] {
	return &decimalField[T, R]{subject: subject}
}

// Type returns the AIP-160 field type for this field.
func (f *decimalField[T, R]) Type() contract.Aip160FieldType {
	return contract.TypeDecimal
}

// ValidateArg validates whether the provided argument is a valid numeric
// decimal.
func (f *decimalField[T, R]) ValidateArg(arg *luciaip160.Arg) error {
	_, err := ormaip160.CoerceArgToNumericConstant(arg)
	return err
}

// BuildCondition builds the query condition for the decimal field.
func (f *decimalField[T, R]) BuildCondition(
	op contract.Aip160Operator,
	arg *luciaip160.Arg,
) (R, error) {
	var zero R
	val, err := ormaip160.CoerceArgToNumericConstant(arg)
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
			"unsupported operator %q for decimal field",
			op,
		)
	}
}
