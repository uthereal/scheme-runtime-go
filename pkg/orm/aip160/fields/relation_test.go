package fields

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"go.chromium.org/luci/common/data/aip160"
)

type dummyField struct {
	typeVal    contract.Aip160FieldType
	validate   error
	buildCond  string
	buildError error
}

func (d *dummyField) Type() contract.Aip160FieldType {
	return d.typeVal
}

func (d *dummyField) ValidateArg(_ *aip160.Arg) error {
	return d.validate
}

func (d *dummyField) BuildCondition(
	_ contract.Aip160Operator,
	_ *aip160.Arg,
) (string, error) {
	return d.buildCond, d.buildError
}

// Test_RelatedField_Type tests the Type method of relatedField.
func Test_RelatedField_Type(t *testing.T) {
	base := &dummyField{typeVal: contract.TypeString}
	f := NewRelatedField[string](base, nil)
	assert.Equal(t, contract.TypeString, f.Type())
}

// Test_RelatedField_ValidateArg tests ValidateArg of relatedField.
func Test_RelatedField_ValidateArg(t *testing.T) {
	base := &dummyField{validate: nil}
	f := NewRelatedField[string](base, nil)
	assert.NoError(t, f.ValidateArg(&aip160.Arg{}))
}

// Test_RelatedField_BuildCondition tests BuildCondition of relatedField.
func Test_RelatedField_BuildCondition(t *testing.T) {
	base := &dummyField{
		buildCond: "child",
	}
	mapper := func(child string) (string, error) {
		return "mapped:" + child, nil
	}
	f := NewRelatedField[string](base, mapper)

	cond, err := f.BuildCondition(contract.OpEq, &aip160.Arg{})
	assert.NoError(t, err)
	assert.Equal(t, "mapped:child", cond)

	// Test error case from base field
	errBase := errors.New("base error")
	errBaseField := &dummyField{
		buildError: errBase,
	}
	fErr := NewRelatedField[string](errBaseField, mapper)
	_, err = fErr.BuildCondition(contract.OpEq, &aip160.Arg{})
	assert.Error(t, err)
}
