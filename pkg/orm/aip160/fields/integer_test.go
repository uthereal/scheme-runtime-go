package fields

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"go.chromium.org/luci/common/data/aip160"
)

type mockIntSubject struct {
	EqFunc  func(int64) string
	NeqFunc func(int64) string
	GtFunc  func(int64) string
	GteFunc func(int64) string
	LtFunc  func(int64) string
	LteFunc func(int64) string
}

func (m mockIntSubject) Eq(val int64) string  { return m.EqFunc(val) }
func (m mockIntSubject) Neq(val int64) string { return m.NeqFunc(val) }
func (m mockIntSubject) Gt(val int64) string  { return m.GtFunc(val) }
func (m mockIntSubject) Gte(val int64) string { return m.GteFunc(val) }
func (m mockIntSubject) Lt(val int64) string  { return m.LtFunc(val) }
func (m mockIntSubject) Lte(val int64) string { return m.LteFunc(val) }

// Test_IntegerField_Type tests the Type method of integerField.
func Test_IntegerField_Type(t *testing.T) {
	subject := mockIntSubject{}
	f := NewIntegerField[int64, string](subject)
	assert.Equal(t, contract.TypeInt, f.Type())
}

// Test_IntegerField_ValidateArg tests ValidateArg on integerField.
func Test_IntegerField_ValidateArg(t *testing.T) {
	subject := mockIntSubject{}
	f := NewIntegerField[int64, string](subject)

	filter, err := aip160.ParseFilter("id = 1234")
	assert.NoError(t, err)
	seq := filter.Expression.Sequences[0]
	arg := seq.Factors[0].Terms[0].Simple.Restriction.Arg

	err = f.ValidateArg(arg)
	assert.NoError(t, err)
}

// Test_IntegerField_BuildCondition tests BuildCondition on integerField.
func Test_IntegerField_BuildCondition(t *testing.T) {
	subject := mockIntSubject{
		EqFunc:  func(v int64) string { return "eq" },
		NeqFunc: func(v int64) string { return "neq" },
		GtFunc:  func(v int64) string { return "gt" },
		GteFunc: func(v int64) string { return "gte" },
		LtFunc:  func(v int64) string { return "lt" },
		LteFunc: func(v int64) string { return "lte" },
	}
	f := NewIntegerField[int64, string](subject)

	filter, err := aip160.ParseFilter("id = 1234")
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
	invalidFilter, err := aip160.ParseFilter("id = \"not-an-int\"")
	assert.NoError(t, err)
	invSeq := invalidFilter.Expression.Sequences[0]
	invalidArg := invSeq.Factors[0].Terms[0].Simple.Restriction.Arg
	_, err = f.BuildCondition(contract.OpEq, invalidArg)
	assert.Error(t, err)
}
