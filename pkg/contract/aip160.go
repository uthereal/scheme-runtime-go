package contract

import (
	"time"
	"uuid"

	"github.com/jackc/pgx/v5/pgtype"
	"go.chromium.org/luci/common/data/aip160"
)

// Aip160Field defines the interface required to build filter conditions.
type Aip160Field[R any] interface {
	// Type returns the AIP-160 field type.
	Type() Aip160FieldType

	// ValidateArg validates the given filter argument.
	ValidateArg(arg *aip160.Arg) error

	// BuildCondition builds the query condition for the operator and argument.
	BuildCondition(op Aip160Operator, arg *aip160.Arg) (R, error)
}

// Aip160StringSubject defines the database column methods required for strings.
type Aip160StringSubject[T ~string, R any] interface {
	// Eq checks if the column is equal to the given value.
	Eq(val T) R

	// Neq checks if the column is not equal to the given value.
	Neq(val T) R
}

// Aip160NumericSubject defines the database column methods required for
// numeric types.
type Aip160NumericSubject[T Number, R any] interface {
	// Eq checks if the column is equal to the given value.
	Eq(val T) R

	// Neq checks if the column is not equal to the given value.
	Neq(val T) R

	// Gt checks if the column is greater than the given value.
	Gt(val T) R

	// Gte checks if the column is greater than or equal to the given value.
	Gte(val T) R

	// Lt checks if the column is less than the given value.
	Lt(val T) R

	// Lte checks if the column is less than or equal to the given value.
	Lte(val T) R
}

// Aip160DecimalSubject defines the database column methods required for exact
// decimal types (pgtype.Numeric).
type Aip160DecimalSubject[T pgtype.Numeric, R any] interface {
	// Eq checks if the column is equal to the given value.
	Eq(val T) R

	// Neq checks if the column is not equal to the given value.
	Neq(val T) R

	// Gt checks if the column is greater than the given value.
	Gt(val T) R

	// Gte checks if the column is greater than or equal to the given value.
	Gte(val T) R

	// Lt checks if the column is less than the given value.
	Lt(val T) R

	// Lte checks if the column is less than or equal to the given value.
	Lte(val T) R
}

// Aip160BoolSubject defines the database column methods required for booleans.
type Aip160BoolSubject[R any] interface {
	// Eq checks if the column is equal to the given value.
	Eq(val bool) R

	// Neq checks if the column is not equal to the given value.
	Neq(val bool) R
}

// Aip160TimestampSubject defines the database column methods required for
// timestamps.
type Aip160TimestampSubject[T time.Time, R any] interface {
	// Eq checks if the column is equal to the given value.
	Eq(val T) R

	// Neq checks if the column is not equal to the given value.
	Neq(val T) R

	// Gt checks if the column is greater than the given value.
	Gt(val T) R

	// Gte checks if the column is greater than or equal to the given value.
	Gte(val T) R

	// Lt checks if the column is less than the given value.
	Lt(val T) R

	// Lte checks if the column is less than or equal to the given value.
	Lte(val T) R
}

// Aip160DurationSubject defines the database column methods required for
// durations.
type Aip160DurationSubject[T time.Duration, R any] interface {
	// Eq checks if the column is equal to the given value.
	Eq(val T) R

	// Neq checks if the column is not equal to the given value.
	Neq(val T) R

	// Gt checks if the column is greater than the given value.
	Gt(val T) R

	// Gte checks if the column is greater than or equal to the given value.
	Gte(val T) R

	// Lt checks if the column is less than the given value.
	Lt(val T) R

	// Lte checks if the column is less than or equal to the given value.
	Lte(val T) R
}

// EnumType represents any string or integer based type used for enum fields.
type EnumType interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Aip160EnumSubject defines the database column methods required for enums.
type Aip160EnumSubject[T any, R any] interface {
	// Eq checks if the column is equal to the given value.
	Eq(val T) R

	// Neq checks if the column is not equal to the given value.
	Neq(val T) R
}

// Aip160UUIDSubject defines the database column methods required for UUID
// fields.
type Aip160UUIDSubject[T uuid.UUID, R any] interface {
	// Eq checks if the column is equal to the given value.
	Eq(val T) R

	// Neq checks if the column is not equal to the given value.
	Neq(val T) R
}

// Aip160FieldType represents the type of filter data.
type Aip160FieldType string

// Aip160Operator represents a comparison operator in a filter expression.
type Aip160Operator string

// TypeBool represents a boolean field type.
const TypeBool Aip160FieldType = "bool"

// TypeDecimal represents a decimal field type.
const TypeDecimal Aip160FieldType = "decimal"

// TypeDuration represents a duration field type.
const TypeDuration Aip160FieldType = "duration"

// TypeEnum represents an enum field type.
const TypeEnum Aip160FieldType = "enum"

// TypeFloat represents a float field type.
const TypeFloat Aip160FieldType = "float"

// TypeInt represents an integer field type.
const TypeInt Aip160FieldType = "int"

// TypeString represents a string field type.
const TypeString Aip160FieldType = "string"

// TypeTimestamp represents a timestamp field type.
const TypeTimestamp Aip160FieldType = "timestamp"

// TypeUUID represents a UUID field type.
const TypeUUID Aip160FieldType = "uuid"

// OpEq represents the equal comparison operator.
const OpEq Aip160Operator = "="

// OpNotEq represents the not-equal comparison operator.
const OpNotEq Aip160Operator = "!="

// OpGt represents the greater-than comparison operator.
const OpGt Aip160Operator = ">"

// OpGtEq represents the greater-than-or-equal comparison operator.
const OpGtEq Aip160Operator = ">="

// OpLt represents the less-than comparison operator.
const OpLt Aip160Operator = "<"

// OpLtEq represents the less-than-or-equal comparison operator.
const OpLtEq Aip160Operator = "<="

// OpHas represents the has/colon operator.
const OpHas Aip160Operator = ":"

// OpDefault represents the default comparison operator.
const OpDefault Aip160Operator = ""
