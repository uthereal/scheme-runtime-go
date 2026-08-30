package aip160

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_HasOperator_Supported tests supported operators for each type.
func Test_HasOperator_Supported(t *testing.T) {
	assert.True(t, HasOperator(contract.TypeBool, contract.OpEq))
	assert.True(t, HasOperator(contract.TypeDecimal, contract.OpGt))
	assert.True(t, HasOperator(contract.TypeDuration, contract.OpLtEq))
	assert.True(t, HasOperator(contract.TypeEnum, contract.OpNotEq))
	assert.True(t, HasOperator(contract.TypeFloat, contract.OpGtEq))
	assert.True(t, HasOperator(contract.TypeInt, contract.OpLt))
	assert.True(t, HasOperator(contract.TypeString, contract.OpHas))
	assert.True(t, HasOperator(contract.TypeTimestamp, contract.OpEq))
	assert.True(t, HasOperator(contract.TypeUUID, contract.OpNotEq))
}

// Test_HasOperator_Unsupported tests unsupported operator scenarios.
func Test_HasOperator_Unsupported(t *testing.T) {
	assert.False(t, HasOperator(contract.TypeBool, contract.OpGt))
	assert.False(t, HasOperator("unknown", contract.OpEq))
	assert.False(t, HasOperator(contract.TypeBool, "unknown"))
}
