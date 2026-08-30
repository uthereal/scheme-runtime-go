package fields

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"go.chromium.org/luci/common/data/aip160"
)

type mockEnumStringSubject struct {
	EqFunc  func(string) string
	NeqFunc func(string) string
}

func (m mockEnumStringSubject) Eq(val string) string {
	return m.EqFunc(val)
}

func (m mockEnumStringSubject) Neq(val string) string {
	return m.NeqFunc(val)
}

// Test_EnumStringField_Type tests the Type method of enumStringField.
func Test_EnumStringField_Type(t *testing.T) {
	subject := mockEnumStringSubject{}
	f := NewStringEnumField[string, string](subject)
	assert.Equal(t, contract.TypeEnum, f.Type())
}

// Test_EnumStringField_ValidateArg tests ValidateArg on enumStringField.
func Test_EnumStringField_ValidateArg(t *testing.T) {
	subject := mockEnumStringSubject{}
	f := NewStringEnumField[string, string](subject)

	filter, err := aip160.ParseFilter("id = \"foo\"")
	assert.NoError(t, err)
	seq := filter.Expression.Sequences[0]
	arg := seq.Factors[0].Terms[0].Simple.Restriction.Arg

	err = f.ValidateArg(arg)
	assert.NoError(t, err)
}

// Test_EnumStringField_BuildCondition tests BuildCondition on enumStringField.
func Test_EnumStringField_BuildCondition(t *testing.T) {
	subject := mockEnumStringSubject{
		EqFunc:  func(v string) string { return "eq" },
		NeqFunc: func(v string) string { return "neq" },
	}
	f := NewStringEnumField[string, string](subject)

	filter, err := aip160.ParseFilter("id = \"foo\"")
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
	invalidArg := &aip160.Arg{}
	_, err = f.BuildCondition(contract.OpEq, invalidArg)
	assert.Error(t, err)
}
