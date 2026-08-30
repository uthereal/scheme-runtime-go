package where

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_InWhere_GetColumn tests GetColumn of InWhere.
func Test_InWhere_GetColumn(t *testing.T) {
	w := InWhere{
		Column: "id",
	}
	assert.Equal(t, "id", w.GetColumn())
}

// Test_InWhere_GetValues tests GetValues of InWhere.
func Test_InWhere_GetValues(t *testing.T) {
	vals := []any{1, 2, 3}
	w := InWhere{
		Values: vals,
	}
	assert.Equal(t, vals, w.GetValues())
}

// Test_InWhere_IsNot tests IsNot of InWhere.
func Test_InWhere_IsNot(t *testing.T) {
	w := InWhere{
		Not: true,
	}
	assert.True(t, w.IsNot())
}

// Test_InWhere_GetBoolean tests GetBoolean of InWhere.
func Test_InWhere_GetBoolean(t *testing.T) {
	w := InWhere{
		Boolean: contract.BoolAnd,
	}
	assert.Equal(t, contract.BoolAnd, w.GetBoolean())
}
