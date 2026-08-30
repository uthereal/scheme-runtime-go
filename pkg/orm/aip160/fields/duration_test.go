package fields

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"go.chromium.org/luci/common/data/aip160"
)

type mockDurationSubject struct {
	EqFunc  func(time.Duration) string
	NeqFunc func(time.Duration) string
	GtFunc  func(time.Duration) string
	GteFunc func(time.Duration) string
	LtFunc  func(time.Duration) string
	LteFunc func(time.Duration) string
}

func (m mockDurationSubject) Eq(val time.Duration) string {
	return m.EqFunc(val)
}

func (m mockDurationSubject) Neq(val time.Duration) string {
	return m.NeqFunc(val)
}

func (m mockDurationSubject) Gt(val time.Duration) string {
	return m.GtFunc(val)
}

func (m mockDurationSubject) Gte(val time.Duration) string {
	return m.GteFunc(val)
}

func (m mockDurationSubject) Lt(val time.Duration) string {
	return m.LtFunc(val)
}

func (m mockDurationSubject) Lte(val time.Duration) string {
	return m.LteFunc(val)
}

// Test_DurationField_Type tests the Type method of durationField.
func Test_DurationField_Type(t *testing.T) {
	subject := mockDurationSubject{}
	f := NewDurationField[time.Duration, string](subject)
	assert.Equal(t, contract.TypeDuration, f.Type())
}

// Test_DurationField_ValidateArg tests ValidateArg on durationField.
func Test_DurationField_ValidateArg(t *testing.T) {
	subject := mockDurationSubject{}
	f := NewDurationField[time.Duration, string](subject)

	filter, err := aip160.ParseFilter("id = \"10s\"")
	assert.NoError(t, err)
	seq := filter.Expression.Sequences[0]
	arg := seq.Factors[0].Terms[0].Simple.Restriction.Arg

	err = f.ValidateArg(arg)
	assert.NoError(t, err)
}

// Test_DurationField_BuildCondition tests BuildCondition on durationField.
func Test_DurationField_BuildCondition(t *testing.T) {
	subject := mockDurationSubject{
		EqFunc:  func(v time.Duration) string { return "eq" },
		NeqFunc: func(v time.Duration) string { return "neq" },
		GtFunc:  func(v time.Duration) string { return "gt" },
		GteFunc: func(v time.Duration) string { return "gte" },
		LtFunc:  func(v time.Duration) string { return "lt" },
		LteFunc: func(v time.Duration) string { return "lte" },
	}
	f := NewDurationField[time.Duration, string](subject)

	filter, err := aip160.ParseFilter("id = \"10s\"")
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
	invalidFilter, err := aip160.ParseFilter("id = \"not-a-duration\"")
	assert.NoError(t, err)
	invSeq := invalidFilter.Expression.Sequences[0]
	invalidArg := invSeq.Factors[0].Terms[0].Simple.Restriction.Arg
	_, err = f.BuildCondition(contract.OpEq, invalidArg)
	assert.Error(t, err)
}
