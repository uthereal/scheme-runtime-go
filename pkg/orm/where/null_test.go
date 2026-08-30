package where

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_NullWhere_GetColumn tests GetColumn of NullWhere.
func Test_NullWhere_GetColumn(t *testing.T) {
	w := NullWhere{
		Column: "deleted_at",
	}
	assert.Equal(t, "deleted_at", w.GetColumn())
}

// Test_NullWhere_IsNot tests IsNot of NullWhere.
func Test_NullWhere_IsNot(t *testing.T) {
	w := NullWhere{
		Not: true,
	}
	assert.True(t, w.IsNot())
}

// Test_NullWhere_GetBoolean tests GetBoolean of NullWhere.
func Test_NullWhere_GetBoolean(t *testing.T) {
	w := NullWhere{
		Boolean: contract.BoolOr,
	}
	assert.Equal(t, contract.BoolOr, w.GetBoolean())
}
