package where

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_BasicWhere_GetColumn tests the GetColumn method of BasicWhere.
func Test_BasicWhere_GetColumn(t *testing.T) {
	w := BasicWhere{
		Column: "id",
	}
	assert.Equal(t, "id", w.GetColumn())
}

// Test_BasicWhere_GetOperator tests the GetOperator method of BasicWhere.
func Test_BasicWhere_GetOperator(t *testing.T) {
	w := BasicWhere{
		Operator: contract.OpEqual,
	}
	assert.Equal(t, contract.OpEqual, w.GetOperator())
}

// Test_BasicWhere_GetValue tests the GetValue method of BasicWhere.
func Test_BasicWhere_GetValue(t *testing.T) {
	w := BasicWhere{
		Value: 42,
	}
	assert.Equal(t, 42, w.GetValue())
}

// Test_BasicWhere_GetBoolean tests the GetBoolean method of BasicWhere.
func Test_BasicWhere_GetBoolean(t *testing.T) {
	w := BasicWhere{
		Boolean: contract.BoolAnd,
	}
	assert.Equal(t, contract.BoolAnd, w.GetBoolean())
}
