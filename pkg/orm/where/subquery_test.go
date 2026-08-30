package where

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_SubQueryWhere_GetColumn tests the GetColumn method.
func Test_SubQueryWhere_GetColumn(t *testing.T) {
	w := SubQueryWhere{
		Column: "category_id",
	}
	assert.Equal(t, "category_id", w.GetColumn())
}

// Test_SubQueryWhere_GetOperator tests the GetOperator method.
func Test_SubQueryWhere_GetOperator(t *testing.T) {
	w := SubQueryWhere{
		Operator: contract.OpEqual,
	}
	assert.Equal(t, contract.OpEqual, w.GetOperator())
}

// Test_SubQueryWhere_GetQuery tests the GetQuery method.
func Test_SubQueryWhere_GetQuery(t *testing.T) {
	type mockQuery struct {
		contract.QueryStateProvider
	}
	m := mockQuery{}
	w := SubQueryWhere{
		Query: m,
	}
	assert.Equal(t, m, w.GetQuery())
}

// Test_SubQueryWhere_GetBoolean tests the GetBoolean method.
func Test_SubQueryWhere_GetBoolean(t *testing.T) {
	w := SubQueryWhere{
		Boolean: contract.BoolAnd,
	}
	assert.Equal(t, contract.BoolAnd, w.GetBoolean())
}
