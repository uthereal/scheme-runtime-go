package fields

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"go.chromium.org/luci/common/data/aip160"
)

type mockStringSubject struct {
	EqFunc  func(string) string
	NeqFunc func(string) string
}

func (m mockStringSubject) Eq(val string) string  { return m.EqFunc(val) }
func (m mockStringSubject) Neq(val string) string { return m.NeqFunc(val) }

// Test_StringField_Type tests the Type method of stringField.
func Test_StringField_Type(t *testing.T) {
	subject := mockStringSubject{}
	f := NewStringField[string, string](subject)
	assert.Equal(t, contract.TypeString, f.Type())
}

// Test_StringField_ValidateArg tests ValidateArg of stringField.
func Test_StringField_ValidateArg(t *testing.T) {
	subject := mockStringSubject{}
	f := NewStringField[string, string](subject)

	filter, err := aip160.ParseFilter("id = \"foo\"")
	assert.NoError(t, err)
	seq := filter.Expression.Sequences[0]
	arg := seq.Factors[0].Terms[0].Simple.Restriction.Arg

	err = f.ValidateArg(arg)
	assert.NoError(t, err)
}

// Test_StringField_BuildCondition tests BuildCondition of stringField.
func Test_StringField_BuildCondition(t *testing.T) {
	subject := mockStringSubject{
		EqFunc:  func(v string) string { return "eq" },
		NeqFunc: func(v string) string { return "neq" },
	}
	f := NewStringField[string, string](subject)

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
