package fields

import (
	"fmt"
	"time"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	ormaip160 "github.com/uthereal/scheme-runtime-go/pkg/orm/aip160"
	luciaip160 "go.chromium.org/luci/common/data/aip160"
)

// durationField is a filter field for duration (time.Duration) columns.
type durationField[T time.Duration, R any] struct {
	subject contract.Aip160DurationSubject[T, R]
}

// NewDurationField creates a new filter field for duration columns.
func NewDurationField[T time.Duration, R any](
	subject contract.Aip160DurationSubject[T, R],
) contract.Aip160Field[R] {
	return &durationField[T, R]{subject: subject}
}

// Type returns the AIP-160 field type for this field.
func (f *durationField[T, R]) Type() contract.Aip160FieldType {
	return contract.TypeDuration
}

// ValidateArg validates whether the provided argument is a valid duration.
func (f *durationField[T, R]) ValidateArg(arg *luciaip160.Arg) error {
	_, err := ormaip160.CoerceArgToDurationConstant(arg)
	return err
}

// BuildCondition builds the query condition for the duration field.
func (f *durationField[T, R]) BuildCondition(
	op contract.Aip160Operator,
	arg *luciaip160.Arg,
) (R, error) {
	var zero R
	val, err := ormaip160.CoerceArgToDurationConstant(arg)
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
			"unsupported operator %q for duration field",
			op,
		)
	}
}
