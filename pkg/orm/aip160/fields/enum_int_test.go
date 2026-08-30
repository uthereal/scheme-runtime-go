package fields

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"go.chromium.org/luci/common/data/aip160"
)

type mockEnumIntSubject struct {
	EqFunc  func(int) string
	NeqFunc func(int) string
}

func (m mockEnumIntSubject) Eq(val int) string {
	return m.EqFunc(val)
}

func (m mockEnumIntSubject) Neq(val int) string {
	return m.NeqFunc(val)
}

// Test_EnumIntField_Type tests the Type method of enumIntField.
func Test_EnumIntField_Type(t *testing.T) {
	subject := mockEnumIntSubject{}
	f := NewIntegerEnumField[int, string](subject)
	assert.Equal(t, contract.TypeEnum, f.Type())
}

// Test_EnumIntField_ValidateArg tests ValidateArg on enumIntField.
func Test_EnumIntField_ValidateArg(t *testing.T) {
	subject := mockEnumIntSubject{}
	f := NewIntegerEnumField[int, string](subject)

	filter, err := aip160.ParseFilter("id = 5")
	assert.NoError(t, err)
	seq := filter.Expression.Sequences[0]
	arg := seq.Factors[0].Terms[0].Simple.Restriction.Arg

	err = f.ValidateArg(arg)
	assert.NoError(t, err)
}

// Test_EnumIntField_BuildCondition tests BuildCondition on enumIntField.
func Test_EnumIntField_BuildCondition(t *testing.T) {
	subject := mockEnumIntSubject{
		EqFunc:  func(v int) string { return "eq" },
		NeqFunc: func(v int) string { return "neq" },
	}
	f := NewIntegerEnumField[int, string](subject)

	filter, err := aip160.ParseFilter("id = 5")
	assert.NoError(t, err)
	seq := filter.Expression.Sequences[0]
	arg := seq.Factors[0].Terms[0].Simple.Restriction.Arg

	cond, err := f.BuildCondition(contract.OpEq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "eq", cond)

	cond, err = f.BuildCondition(contract.OpNotEq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "neq", cond)

	// Test unsupported operator
	_, err = f.BuildCondition(contract.OpGt, arg)
	assert.Error(t, err)

	// Test coercion error
	invalidFilter, err := aip160.ParseFilter("id = \"not-an-int\"")
	assert.NoError(t, err)
	invSeq := invalidFilter.Expression.Sequences[0]
	invalidArg := invSeq.Factors[0].Terms[0].Simple.Restriction.Arg
	_, err = f.BuildCondition(contract.OpEq, invalidArg)
	assert.Error(t, err)
}
