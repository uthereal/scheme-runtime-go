package fields

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"go.chromium.org/luci/common/data/aip160"
)

type mockUUIDSubject struct {
	EqFunc  func(uuid.UUID) string
	NeqFunc func(uuid.UUID) string
}

func (m mockUUIDSubject) Eq(val uuid.UUID) string  { return m.EqFunc(val) }
func (m mockUUIDSubject) Neq(val uuid.UUID) string { return m.NeqFunc(val) }

// Test_UUIDField_Type tests the Type method of uuidField.
func Test_UUIDField_Type(t *testing.T) {
	subject := mockUUIDSubject{}
	f := NewUUIDField[uuid.UUID, string](subject)
	assert.Equal(t, contract.TypeUUID, f.Type())
}

// Test_UUIDField_ValidateArg tests ValidateArg on uuidField.
func Test_UUIDField_ValidateArg(t *testing.T) {
	subject := mockUUIDSubject{}
	f := NewUUIDField[uuid.UUID, string](subject)

	filter, err := aip160.ParseFilter(
		"id = \"550e8400-e29b-41d4-a716-446655440000\"",
	)
	assert.NoError(t, err)
	seq := filter.Expression.Sequences[0]
	arg := seq.Factors[0].Terms[0].Simple.Restriction.Arg

	err = f.ValidateArg(arg)
	assert.NoError(t, err)

	// Test coercion error
	err = f.ValidateArg(&aip160.Arg{})
	assert.Error(t, err)
}

// Test_UUIDField_BuildCondition tests BuildCondition on uuidField.
func Test_UUIDField_BuildCondition(t *testing.T) {
	subject := mockUUIDSubject{
		EqFunc:  func(v uuid.UUID) string { return "eq" },
		NeqFunc: func(v uuid.UUID) string { return "neq" },
	}
	f := NewUUIDField[uuid.UUID, string](subject)

	filter, err := aip160.ParseFilter(
		"id = \"550e8400-e29b-41d4-a716-446655440000\"",
	)
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

	// Test parse UUID error
	invalidFilter, err := aip160.ParseFilter("id = \"not-a-uuid\"")
	assert.NoError(t, err)
	invSeq := invalidFilter.Expression.Sequences[0]
	invalidArg := invSeq.Factors[0].Terms[0].Simple.Restriction.Arg
	_, err = f.BuildCondition(contract.OpEq, invalidArg)
	assert.Error(t, err)

	// Test coercion error
	_, err = f.BuildCondition(contract.OpEq, &aip160.Arg{})
	assert.Error(t, err)
}
