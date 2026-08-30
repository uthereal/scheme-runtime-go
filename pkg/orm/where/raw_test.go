package where

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_RawWhere_GetSql tests the GetSql method of RawWhere.
func Test_RawWhere_GetSql(t *testing.T) {
	w := RawWhere{
		Sql: "id = ?",
	}
	assert.Equal(t, "id = ?", w.GetSql())
}

// Test_RawWhere_GetBindings tests the GetBindings method of RawWhere.
func Test_RawWhere_GetBindings(t *testing.T) {
	bindings := []any{42}
	w := RawWhere{
		Bindings: bindings,
	}
	assert.Equal(t, bindings, w.GetBindings())
}

// Test_RawWhere_GetBoolean tests the GetBoolean method of RawWhere.
func Test_RawWhere_GetBoolean(t *testing.T) {
	w := RawWhere{
		Boolean: contract.BoolAnd,
	}
	assert.Equal(t, contract.BoolAnd, w.GetBoolean())
}
