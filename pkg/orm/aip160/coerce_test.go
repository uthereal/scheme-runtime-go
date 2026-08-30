package aip160

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.chromium.org/luci/common/data/aip160"
)

// Test_Coerce_StringConstant tests coercion functions on String types.
func Test_Coerce_StringConstant(t *testing.T) {
	f, err := aip160.ParseFilter("id = \"foo\"")
	assert.NoError(t, err)
	arg := f.Expression.Sequences[0].Factors[0].Terms[0].Simple.Restriction.Arg

	val, err := CoerceArgToStringConstant(arg)
	assert.NoError(t, err)
	assert.Equal(t, "foo", val)
}

// Test_Coerce_FloatConstant tests coercion functions on Float types.
func Test_Coerce_FloatConstant(t *testing.T) {
	f, err := aip160.ParseFilter("id = 12.34")
	assert.NoError(t, err)
	arg := f.Expression.Sequences[0].Factors[0].Terms[0].Simple.Restriction.Arg

	val, err := CoerceArgToFloatConstant(arg)
	assert.NoError(t, err)
	assert.Equal(t, 12.34, val)
}

// Test_Coerce_IntConstant tests coercion functions on Int types.
func Test_Coerce_IntConstant(t *testing.T) {
	f, err := aip160.ParseFilter("id = 1234")
	assert.NoError(t, err)
	arg := f.Expression.Sequences[0].Factors[0].Terms[0].Simple.Restriction.Arg

	val, err := CoerceArgToIntConstant(arg)
	assert.NoError(t, err)
	assert.Equal(t, int64(1234), val)
}

// Test_Coerce_BoolConstant tests coercion functions on Bool types.
func Test_Coerce_BoolConstant(t *testing.T) {
	f, err := aip160.ParseFilter("id = true")
	assert.NoError(t, err)
	arg := f.Expression.Sequences[0].Factors[0].Terms[0].Simple.Restriction.Arg

	val, err := CoerceArgToBoolConstant(arg)
	assert.NoError(t, err)
	assert.True(t, val)

	f2, err := aip160.ParseFilter("id = false")
	assert.NoError(t, err)
	arg2 := f2.Expression.Sequences[0].Factors[0].Terms[0].Simple.Restriction.Arg

	val2, err := CoerceArgToBoolConstant(arg2)
	assert.NoError(t, err)
	assert.False(t, val2)
}

// Test_Coerce_TimestampConstant tests coercion functions on Timestamp types.
func Test_Coerce_TimestampConstant(t *testing.T) {
	f, err := aip160.ParseFilter("id = \"2026-08-08T12:34:56Z\"")
	assert.NoError(t, err)
	arg := f.Expression.Sequences[0].Factors[0].Terms[0].Simple.Restriction.Arg

	val, err := CoerceArgToTimestampConstant(arg)
	assert.NoError(t, err)
	assert.Equal(t, int64(2026), int64(val.Year()))
}

// Test_Coerce_DurationConstant tests coercion functions on Duration types.
func Test_Coerce_DurationConstant(t *testing.T) {
	f, err := aip160.ParseFilter("id = \"10s\"")
	assert.NoError(t, err)
	arg := f.Expression.Sequences[0].Factors[0].Terms[0].Simple.Restriction.Arg

	val, err := CoerceArgToDurationConstant(arg)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), int64(val.Seconds()))

	f2, err := aip160.ParseFilter("id = 30")
	assert.NoError(t, err)
	arg2 := f2.Expression.Sequences[0].Factors[0].Terms[0].Simple.Restriction.Arg

	val2, err := CoerceArgToDurationConstant(arg2)
	assert.NoError(t, err)
	assert.Equal(t, int64(30), int64(val2.Seconds()))
}

// Test_Coerce_NumericConstant tests coercion functions on Numeric decimal
// types.
func Test_Coerce_NumericConstant(t *testing.T) {
	f, err := aip160.ParseFilter("id = 12.3456")
	assert.NoError(t, err)
	arg := f.Expression.Sequences[0].Factors[0].Terms[0].Simple.Restriction.Arg

	val, err := CoerceArgToNumericConstant(arg)
	assert.NoError(t, err)
	assert.NotNil(t, val)
}

// Test_Coerce_Errors tests invalid coercion scenario errors.
func Test_Coerce_Errors(t *testing.T) {
	// 1. CoerceArgToComparable errors
	compositeArg := &aip160.Arg{
		Composite: &aip160.Expression{},
	}
	_, err := CoerceArgToComparable(compositeArg)
	assert.Error(t, err)

	missingCmpArg := &aip160.Arg{}
	_, err = CoerceArgToComparable(missingCmpArg)
	assert.Error(t, err)

	// Helper function to create an Arg with a Comparable having nil Member
	argWithNilMember := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: nil,
		},
	}

	// Helper function to create an Arg with a Comparable having nil Value
	argWithNilValue := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: nil,
			},
		},
	}

	// 2. CoerceArgToStringConstant errors
	_, err = CoerceArgToStringConstant(missingCmpArg)
	assert.Error(t, err)

	_, err = CoerceArgToStringConstant(argWithNilMember)
	assert.Error(t, err)

	_, err = CoerceArgToStringConstant(argWithNilValue)
	assert.Error(t, err)

	// 3. CoerceArgToFloatConstant errors
	_, err = CoerceArgToFloatConstant(missingCmpArg)
	assert.Error(t, err)

	_, err = CoerceArgToFloatConstant(argWithNilMember)
	assert.Error(t, err)

	_, err = CoerceArgToFloatConstant(argWithNilValue)
	assert.Error(t, err)

	quotedFloatArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "12.34",
					Quoted: true,
				},
			},
		},
	}
	_, err = CoerceArgToFloatConstant(quotedFloatArg)
	assert.Error(t, err)

	emptyFloatArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "",
					Quoted: false,
				},
			},
		},
	}
	_, err = CoerceArgToFloatConstant(emptyFloatArg)
	assert.Error(t, err)

	invalidFloatArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "12.34.56",
					Quoted: false,
				},
			},
		},
	}
	_, err = CoerceArgToFloatConstant(invalidFloatArg)
	assert.Error(t, err)

	// 4. CoerceArgToIntConstant errors
	_, err = CoerceArgToIntConstant(missingCmpArg)
	assert.Error(t, err)

	_, err = CoerceArgToIntConstant(argWithNilMember)
	assert.Error(t, err)

	_, err = CoerceArgToIntConstant(argWithNilValue)
	assert.Error(t, err)

	quotedIntArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "123",
					Quoted: true,
				},
			},
		},
	}
	_, err = CoerceArgToIntConstant(quotedIntArg)
	assert.Error(t, err)

	emptyIntArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "",
					Quoted: false,
				},
			},
		},
	}
	_, err = CoerceArgToIntConstant(emptyIntArg)
	assert.Error(t, err)

	invalidIntArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "12a",
					Quoted: false,
				},
			},
		},
	}
	_, err = CoerceArgToIntConstant(invalidIntArg)
	assert.Error(t, err)

	// 5. CoerceArgToBoolConstant errors
	_, err = CoerceArgToBoolConstant(missingCmpArg)
	assert.Error(t, err)

	_, err = CoerceArgToBoolConstant(argWithNilMember)
	assert.Error(t, err)

	invalidBoolArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "not-a-bool",
					Quoted: false,
				},
			},
		},
	}
	_, err = CoerceArgToBoolConstant(invalidBoolArg)
	assert.Error(t, err)

	// 6. CoerceArgToTimestampConstant errors
	_, err = CoerceArgToTimestampConstant(missingCmpArg)
	assert.Error(t, err)

	_, err = CoerceArgToTimestampConstant(argWithNilMember)
	assert.Error(t, err)

	_, err = CoerceArgToTimestampConstant(argWithNilValue)
	assert.Error(t, err)

	unquotedTimestampArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "2026-08-08T12:34:56Z",
					Quoted: false,
				},
			},
		},
	}
	_, err = CoerceArgToTimestampConstant(unquotedTimestampArg)
	assert.Error(t, err)

	emptyTimestampArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "",
					Quoted: true,
				},
			},
		},
	}
	_, err = CoerceArgToTimestampConstant(emptyTimestampArg)
	assert.Error(t, err)

	invalidTimestampArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "2026-08-08",
					Quoted: true,
				},
			},
		},
	}
	_, err = CoerceArgToTimestampConstant(invalidTimestampArg)
	assert.Error(t, err)

	// 7. CoerceArgToDurationConstant errors
	_, err = CoerceArgToDurationConstant(missingCmpArg)
	assert.Error(t, err)

	invalidDurationCmpArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: nil,
		},
	}
	_, err = CoerceArgToDurationConstant(invalidDurationCmpArg)
	assert.Error(t, err)

	_, err = CoerceArgToDurationConstant(argWithNilValue)
	assert.Error(t, err)

	quotedInvalidDurationArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "10xx",
					Quoted: true,
				},
			},
		},
	}
	_, err = CoerceArgToDurationConstant(quotedInvalidDurationArg)
	assert.Error(t, err)

	unquotedFloatDurationArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "1.5",
					Quoted: false,
				},
			},
		},
	}
	_, err = CoerceArgToDurationConstant(unquotedFloatDurationArg)
	assert.NoError(t, err)

	unquotedInvalidDurationArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "foo",
					Quoted: false,
				},
			},
		},
	}
	_, err = CoerceArgToDurationConstant(unquotedInvalidDurationArg)
	assert.Error(t, err)

	// 8. CoerceArgToNumericConstant errors
	_, err = CoerceArgToNumericConstant(missingCmpArg)
	assert.Error(t, err)

	_, err = CoerceArgToNumericConstant(argWithNilMember)
	assert.Error(t, err)

	_, err = CoerceArgToNumericConstant(argWithNilValue)
	assert.Error(t, err)

	quotedNumericArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "12.34",
					Quoted: true,
				},
			},
		},
	}
	_, err = CoerceArgToNumericConstant(quotedNumericArg)
	assert.Error(t, err)

	emptyNumericArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "",
					Quoted: false,
				},
			},
		},
	}
	_, err = CoerceArgToNumericConstant(emptyNumericArg)
	assert.Error(t, err)

	invalidNumericArg := &aip160.Arg{
		Comparable: &aip160.Comparable{
			Member: &aip160.Member{
				Value: &aip160.Value{
					Value:  "invalid-dec",
					Quoted: false,
				},
			},
		},
	}
	_, err = CoerceArgToNumericConstant(invalidNumericArg)
	assert.Error(t, err)
}
