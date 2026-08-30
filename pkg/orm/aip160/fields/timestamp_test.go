package fields

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"go.chromium.org/luci/common/data/aip160"
)

type mockTimestampSubject struct {
	EqFunc  func(time.Time) string
	NeqFunc func(time.Time) string
	GtFunc  func(time.Time) string
	GteFunc func(time.Time) string
	LtFunc  func(time.Time) string
	LteFunc func(time.Time) string
}

func (m mockTimestampSubject) Eq(val time.Time) string  { return m.EqFunc(val) }
func (m mockTimestampSubject) Neq(val time.Time) string {
	return m.NeqFunc(val)
}
func (m mockTimestampSubject) Gt(val time.Time) string  { return m.GtFunc(val) }
func (m mockTimestampSubject) Gte(val time.Time) string {
	return m.GteFunc(val)
}
func (m mockTimestampSubject) Lt(val time.Time) string  { return m.LtFunc(val) }
func (m mockTimestampSubject) Lte(val time.Time) string {
	return m.LteFunc(val)
}

// Test_TimestampField_Type tests the Type method of timestampField.
func Test_TimestampField_Type(t *testing.T) {
	subject := mockTimestampSubject{}
	f := NewTimestampField[time.Time, string](subject)
	assert.Equal(t, contract.TypeTimestamp, f.Type())
}

// Test_TimestampField_ValidateArg tests ValidateArg on timestampField.
func Test_TimestampField_ValidateArg(t *testing.T) {
	subject := mockTimestampSubject{}
	f := NewTimestampField[time.Time, string](subject)

	filter, err := aip160.ParseFilter("id = \"2026-08-08T12:34:56Z\"")
	assert.NoError(t, err)
	seq := filter.Expression.Sequences[0]
	arg := seq.Factors[0].Terms[0].Simple.Restriction.Arg

	err = f.ValidateArg(arg)
	assert.NoError(t, err)
}

// Test_TimestampField_BuildCondition tests BuildCondition on timestampField.
func Test_TimestampField_BuildCondition(t *testing.T) {
	subject := mockTimestampSubject{
		EqFunc:  func(v time.Time) string { return "eq" },
		NeqFunc: func(v time.Time) string { return "neq" },
		GtFunc:  func(v time.Time) string { return "gt" },
		GteFunc: func(v time.Time) string { return "gte" },
		LtFunc:  func(v time.Time) string { return "lt" },
		LteFunc: func(v time.Time) string { return "lte" },
	}
	f := NewTimestampField[time.Time, string](subject)

	filter, err := aip160.ParseFilter("id = \"2026-08-08T12:34:56Z\"")
	assert.NoError(t, err)
	seq := filter.Expression.Sequences[0]
	arg := seq.Factors[0].Terms[0].Simple.Restriction.Arg

	for _, op := range []contract.Aip160Operator{
		contract.OpEq, contract.OpNotEq, contract.OpGt,
		contract.OpGtEq, contract.OpLt, contract.OpLtEq,
	} {
		cond, err := f.BuildCondition(op, arg)
		assert.NoError(t, err)
		assert.NotEmpty(t, cond)
	}

	// Test unsupported operator
	_, err = f.BuildCondition(contract.OpHas, arg)
	assert.Error(t, err)

	// Test coercion error
	invalidFilter, err := aip160.ParseFilter("id = \"not-a-timestamp\"")
	assert.NoError(t, err)
	invSeq := invalidFilter.Expression.Sequences[0]
	invalidArg := invSeq.Factors[0].Terms[0].Simple.Restriction.Arg
	_, err = f.BuildCondition(contract.OpEq, invalidArg)
	assert.Error(t, err)
}
