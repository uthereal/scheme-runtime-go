package fields

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"go.chromium.org/luci/common/data/aip160"
)

type mockBoolSubject struct {
	EqFunc  func(bool) string
	NeqFunc func(bool) string
}

func (m mockBoolSubject) Eq(val bool) string {
	return m.EqFunc(val)
}

func (m mockBoolSubject) Neq(val bool) string {
	return m.NeqFunc(val)
}

// Test_BoolField_Type tests the Type method of boolField.
func Test_BoolField_Type(t *testing.T) {
	subject := mockBoolSubject{}
	f := NewBoolField[string](subject)
	assert.Equal(t, contract.TypeBool, f.Type())
}

// Test_BoolField_ValidateArg tests ValidateArg on boolField.
func Test_BoolField_ValidateArg(t *testing.T) {
	subject := mockBoolSubject{}
	f := NewBoolField[string](subject)

	filter, err := aip160.ParseFilter("id = true")
	assert.NoError(t, err)
	seq := filter.Expression.Sequences[0]
	arg := seq.Factors[0].Terms[0].Simple.Restriction.Arg

	err = f.ValidateArg(arg)
	assert.NoError(t, err)

	invalidFilter, err := aip160.ParseFilter("id = \"not-a-bool\"")
	assert.NoError(t, err)
	invSeq := invalidFilter.Expression.Sequences[0]
	invalidArg := invSeq.Factors[0].Terms[0].Simple.Restriction.Arg

	err = f.ValidateArg(invalidArg)
	assert.Error(t, err)
}

// Test_BoolField_BuildCondition tests BuildCondition on boolField.
func Test_BoolField_BuildCondition(t *testing.T) {
	subject := mockBoolSubject{
		EqFunc: func(v bool) string {
			return fmt.Sprintf("eq(%t)", v)
		},
		NeqFunc: func(v bool) string {
			return fmt.Sprintf("neq(%t)", v)
		},
	}
	f := NewBoolField[string](subject)

	filter, err := aip160.ParseFilter("id = true")
	assert.NoError(t, err)
	seq := filter.Expression.Sequences[0]
	arg := seq.Factors[0].Terms[0].Simple.Restriction.Arg

	// Test Eq
	cond, err := f.BuildCondition(contract.OpEq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "eq(true)", cond)

	// Test Neq
	cond, err = f.BuildCondition(contract.OpNotEq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "neq(true)", cond)

	// Test Unsupported Operator
	_, err = f.BuildCondition(contract.OpGt, arg)
	assert.Error(t, err)

	// Test Coercion Error
	invalidFilter, err := aip160.ParseFilter("id = \"not-a-bool\"")
	assert.NoError(t, err)
	invSeq := invalidFilter.Expression.Sequences[0]
	invalidArg := invSeq.Factors[0].Terms[0].Simple.Restriction.Arg
	_, err = f.BuildCondition(contract.OpEq, invalidArg)
	assert.Error(t, err)
}
