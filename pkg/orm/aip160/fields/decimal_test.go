package fields

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"go.chromium.org/luci/common/data/aip160"
)

type mockDecimalSubject struct {
	EqFunc  func(pgtype.Numeric) string
	NeqFunc func(pgtype.Numeric) string
	GtFunc  func(pgtype.Numeric) string
	GteFunc func(pgtype.Numeric) string
	LtFunc  func(pgtype.Numeric) string
	LteFunc func(pgtype.Numeric) string
}

func (m mockDecimalSubject) Eq(val pgtype.Numeric) string {
	return m.EqFunc(val)
}

func (m mockDecimalSubject) Neq(val pgtype.Numeric) string {
	return m.NeqFunc(val)
}

func (m mockDecimalSubject) Gt(val pgtype.Numeric) string {
	return m.GtFunc(val)
}

func (m mockDecimalSubject) Gte(val pgtype.Numeric) string {
	return m.GteFunc(val)
}

func (m mockDecimalSubject) Lt(val pgtype.Numeric) string {
	return m.LtFunc(val)
}

func (m mockDecimalSubject) Lte(val pgtype.Numeric) string {
	return m.LteFunc(val)
}

// Test_DecimalField_Type tests the Type method of decimalField.
func Test_DecimalField_Type(t *testing.T) {
	subject := mockDecimalSubject{}
	f := NewDecimalField[pgtype.Numeric, string](subject)
	assert.Equal(t, contract.TypeDecimal, f.Type())
}

// Test_DecimalField_ValidateArg tests ValidateArg on decimalField.
func Test_DecimalField_ValidateArg(t *testing.T) {
	subject := mockDecimalSubject{}
	f := NewDecimalField[pgtype.Numeric, string](subject)

	filter, err := aip160.ParseFilter("id = 12.34")
	assert.NoError(t, err)
	seq := filter.Expression.Sequences[0]
	arg := seq.Factors[0].Terms[0].Simple.Restriction.Arg

	err = f.ValidateArg(arg)
	assert.NoError(t, err)
}

// Test_DecimalField_BuildCondition tests BuildCondition on decimalField.
func Test_DecimalField_BuildCondition(t *testing.T) {
	subject := mockDecimalSubject{
		EqFunc:  func(v pgtype.Numeric) string { return "eq" },
		NeqFunc: func(v pgtype.Numeric) string { return "neq" },
		GtFunc:  func(v pgtype.Numeric) string { return "gt" },
		GteFunc: func(v pgtype.Numeric) string { return "gte" },
		LtFunc:  func(v pgtype.Numeric) string { return "lt" },
		LteFunc: func(v pgtype.Numeric) string { return "lte" },
	}
	f := NewDecimalField[pgtype.Numeric, string](subject)

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
	invalidFilter, err := aip160.ParseFilter("id = \"not-a-decimal\"")
	assert.NoError(t, err)
	invSeq := invalidFilter.Expression.Sequences[0]
	invalidArg := invSeq.Factors[0].Terms[0].Simple.Restriction.Arg
	_, err = f.BuildCondition(contract.OpEq, invalidArg)
	assert.Error(t, err)
}
