package fields

import (
	"fmt"
	"time"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	ormaip160 "github.com/uthereal/scheme-runtime-go/pkg/orm/aip160"
	luciaip160 "go.chromium.org/luci/common/data/aip160"
)

// timestampField represents a filter field for timestamp columns.
type timestampField[T time.Time, R any] struct {
	subject contract.Aip160TimestampSubject[T, R]
}

// NewTimestampField creates a new filter field for timestamp (time.Time)
// columns.
func NewTimestampField[
	T time.Time,
	R any,
](
	subject contract.Aip160TimestampSubject[T, R],
) contract.Aip160Field[R] {
	return &timestampField[T, R]{subject: subject}
}

// Type returns the field type.
func (f *timestampField[T, R]) Type() contract.Aip160FieldType {
	return contract.TypeTimestamp
}

// ValidateArg validates that the provided argument can be coerced to a
// timestamp.
func (f *timestampField[T, R]) ValidateArg(arg *luciaip160.Arg) error {
	_, err := ormaip160.CoerceArgToTimestampConstant(arg)
	return err
}

// BuildCondition builds an AIP-160 filter condition for a timestamp field.
func (f *timestampField[T, R]) BuildCondition(
	op contract.Aip160Operator,
	arg *luciaip160.Arg,
) (R, error) {
	var zero R
	val, err := ormaip160.CoerceArgToTimestampConstant(arg)
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
			"unsupported operator %q for timestamp field",
			op,
		)
	}
}
