package where

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_DateWhere_GetType tests the GetType method of DateWhere.
func Test_DateWhere_GetType(t *testing.T) {
	w := DateWhere{
		Type: "year",
	}
	assert.Equal(t, "year", w.GetType())
}

// Test_DateWhere_GetColumn tests the GetColumn method of DateWhere.
func Test_DateWhere_GetColumn(t *testing.T) {
	w := DateWhere{
		Column: "created_at",
	}
	assert.Equal(t, "created_at", w.GetColumn())
}

// Test_DateWhere_GetOperator tests the GetOperator method of DateWhere.
func Test_DateWhere_GetOperator(t *testing.T) {
	w := DateWhere{
		Operator: contract.OpEqual,
	}
	assert.Equal(t, contract.OpEqual, w.GetOperator())
}

// Test_DateWhere_GetValue tests the GetValue method of DateWhere.
func Test_DateWhere_GetValue(t *testing.T) {
	w := DateWhere{
		Value: "2023",
	}
	assert.Equal(t, "2023", w.GetValue())
}

// Test_DateWhere_GetBoolean tests the GetBoolean method of DateWhere.
func Test_DateWhere_GetBoolean(t *testing.T) {
	w := DateWhere{
		Boolean: contract.BoolAnd,
	}
	assert.Equal(t, contract.BoolAnd, w.GetBoolean())
}
