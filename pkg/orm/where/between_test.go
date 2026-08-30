package where

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_BetweenWhere_GetColumn tests the GetColumn method of BetweenWhere.
func Test_BetweenWhere_GetColumn(t *testing.T) {
	w := BetweenWhere{
		Column: "age",
	}
	assert.Equal(t, "age", w.GetColumn())
}

// Test_BetweenWhere_GetMin tests the GetMin method of BetweenWhere.
func Test_BetweenWhere_GetMin(t *testing.T) {
	w := BetweenWhere{
		Min: 18,
	}
	assert.Equal(t, 18, w.GetMin())
}

// Test_BetweenWhere_GetMax tests the GetMax method of BetweenWhere.
func Test_BetweenWhere_GetMax(t *testing.T) {
	w := BetweenWhere{
		Max: 65,
	}
	assert.Equal(t, 65, w.GetMax())
}

// Test_BetweenWhere_IsNot tests the IsNot method of BetweenWhere.
func Test_BetweenWhere_IsNot(t *testing.T) {
	w := BetweenWhere{
		Not: true,
	}
	assert.True(t, w.IsNot())
}

// Test_BetweenWhere_GetBoolean tests the GetBoolean method of BetweenWhere.
func Test_BetweenWhere_GetBoolean(t *testing.T) {
	w := BetweenWhere{
		Boolean: contract.BoolAnd,
	}
	assert.Equal(t, contract.BoolAnd, w.GetBoolean())
}
