package fields

import (
	"fmt"
	"uuid"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	ormaip160 "github.com/uthereal/scheme-runtime-go/pkg/orm/aip160"
	luciaip160 "go.chromium.org/luci/common/data/aip160"
)

// uuidField represents a filter field for UUID columns.
type uuidField[T uuid.UUID, R any] struct {
	subject contract.Aip160UUIDSubject[T, R]
}

// NewUUIDField creates a new filter field for UUID columns.
func NewUUIDField[
	T uuid.UUID,
	R any,
](
	subject contract.Aip160UUIDSubject[T, R],
) contract.Aip160Field[R] {
	return &uuidField[T, R]{subject: subject}
}

// Type returns the field type.
func (f *uuidField[T, R]) Type() contract.Aip160FieldType {
	return contract.TypeUUID
}

// ValidateArg validates that the provided argument can be parsed as a
// UUID.
func (f *uuidField[T, R]) ValidateArg(arg *luciaip160.Arg) error {
	val, err := ormaip160.CoerceArgToStringConstant(arg)
	if err != nil {
		return err
	}
	_, err = uuid.Parse(val)
	return err
}

// BuildCondition builds an AIP-160 filter condition for a UUID field.
func (f *uuidField[T, R]) BuildCondition(
	op contract.Aip160Operator,
	arg *luciaip160.Arg,
) (R, error) {
	var zero R
	val, err := ormaip160.CoerceArgToStringConstant(arg)
	if err != nil {
		return zero, err
	}
	parsed, err := uuid.Parse(val)
	if err != nil {
		return zero, fmt.Errorf(
			"failed to parse UUID value %q -> %w",
			val,
			err,
		)
	}
	switch op {
	case contract.OpEq, contract.OpDefault:
		return f.subject.Eq(T(parsed)), nil
	case contract.OpNotEq:
		return f.subject.Neq(T(parsed)), nil
	default:
		return zero, fmt.Errorf(
			"unsupported operator %q for uuid field",
			op,
		)
	}
}
