package where

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_ColumnWhere_GetFirst tests the GetFirst method of ColumnWhere.
func Test_ColumnWhere_GetFirst(t *testing.T) {
	w := ColumnWhere{
		First: "created_at",
	}
	assert.Equal(t, "created_at", w.GetFirst())
}

// Test_ColumnWhere_GetOperator tests the GetOperator method of ColumnWhere.
func Test_ColumnWhere_GetOperator(t *testing.T) {
	w := ColumnWhere{
		Operator: contract.OpEqual,
	}
	assert.Equal(t, contract.OpEqual, w.GetOperator())
}

// Test_ColumnWhere_GetSecond tests the GetSecond method of ColumnWhere.
func Test_ColumnWhere_GetSecond(t *testing.T) {
	w := ColumnWhere{
		Second: "updated_at",
	}
	assert.Equal(t, "updated_at", w.GetSecond())
}

// Test_ColumnWhere_GetBoolean tests the GetBoolean method of ColumnWhere.
func Test_ColumnWhere_GetBoolean(t *testing.T) {
	w := ColumnWhere{
		Boolean: contract.BoolAnd,
	}
	assert.Equal(t, contract.BoolAnd, w.GetBoolean())
}
