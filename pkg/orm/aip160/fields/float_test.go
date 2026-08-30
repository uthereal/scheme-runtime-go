package fields

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"go.chromium.org/luci/common/data/aip160"
)

type mockFloatSubject struct {
	EqFunc  func(float64) string
	NeqFunc func(float64) string
	GtFunc  func(float64) string
	GteFunc func(float64) string
	LtFunc  func(float64) string
	LteFunc func(float64) string
}

func (m mockFloatSubject) Eq(val float64) string  { return m.EqFunc(val) }
func (m mockFloatSubject) Neq(val float64) string { return m.NeqFunc(val) }
func (m mockFloatSubject) Gt(val float64) string  { return m.GtFunc(val) }
func (m mockFloatSubject) Gte(val float64) string { return m.GteFunc(val) }
func (m mockFloatSubject) Lt(val float64) string  { return m.LtFunc(val) }
func (m mockFloatSubject) Lte(val float64) string { return m.LteFunc(val) }

// Test_FloatField_Type tests the Type method of floatField.
func Test_FloatField_Type(t *testing.T) {
	subject := mockFloatSubject{}
	f := NewFloatField[float64, string](subject)
	assert.Equal(t, contract.TypeFloat, f.Type())
}

// Test_FloatField_ValidateArg tests ValidateArg on floatField.
func Test_FloatField_ValidateArg(t *testing.T) {
	subject := mockFloatSubject{}
	f := NewFloatField[float64, string](subject)

	filter, err := aip160.ParseFilter("id = 12.34")
	assert.NoError(t, err)
	seq := filter.Expression.Sequences[0]
	arg := seq.Factors[0].Terms[0].Simple.Restriction.Arg

	err = f.ValidateArg(arg)
	assert.NoError(t, err)
}

// Test_FloatField_BuildCondition tests BuildCondition on floatField.
func Test_FloatField_BuildCondition(t *testing.T) {
	subject := mockFloatSubject{
		EqFunc:  func(v float64) string { return "eq" },
		NeqFunc: func(v float64) string { return "neq" },
		GtFunc:  func(v float64) string { return "gt" },
		GteFunc: func(v float64) string { return "gte" },
		LtFunc:  func(v float64) string { return "lt" },
		LteFunc: func(v float64) string { return "lte" },
	}
	f := NewFloatField[float64, string](subject)

	filter, err := aip160.ParseFilter("id = 12.34")
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
	invalidFilter, err := aip160.ParseFilter("id = \"not-a-float\"")
	assert.NoError(t, err)
	invSeq := invalidFilter.Expression.Sequences[0]
	invalidArg := invSeq.Factors[0].Terms[0].Simple.Restriction.Arg
	_, err = f.BuildCondition(contract.OpEq, invalidArg)
	assert.Error(t, err)
}
