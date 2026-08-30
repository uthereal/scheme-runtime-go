package where

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_JsonWhere_GetColumn tests GetColumn of JsonWhere.
func Test_JsonWhere_GetColumn(t *testing.T) {
	w := JsonWhere{
		Column: "data",
	}
	assert.Equal(t, "data", w.GetColumn())
}

// Test_JsonWhere_GetKey tests GetKey of JsonWhere.
func Test_JsonWhere_GetKey(t *testing.T) {
	w := JsonWhere{
		Key: "status",
	}
	assert.Equal(t, "status", w.GetKey())
}

// Test_JsonWhere_GetOperator tests GetOperator of JsonWhere.
func Test_JsonWhere_GetOperator(t *testing.T) {
	w := JsonWhere{
		Operator: contract.OpEqual,
	}
	assert.Equal(t, contract.OpEqual, w.GetOperator())
}

// Test_JsonWhere_GetValue tests GetValue of JsonWhere.
func Test_JsonWhere_GetValue(t *testing.T) {
	w := JsonWhere{
		Value: "active",
	}
	assert.Equal(t, "active", w.GetValue())
}

// Test_JsonWhere_GetBoolean tests GetBoolean of JsonWhere.
func Test_JsonWhere_GetBoolean(t *testing.T) {
	w := JsonWhere{
		Boolean: contract.BoolOr,
	}
	assert.Equal(t, contract.BoolOr, w.GetBoolean())
}
