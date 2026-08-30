package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_ColumnOrder_GetColumn tests the GetColumn method of ColumnOrder.
func Test_ColumnOrder_GetColumn(t *testing.T) {
	o := ColumnOrder{
		Column: "created_at",
	}
	assert.Equal(t, "created_at", o.GetColumn())
}

// Test_ColumnOrder_GetDirection tests the GetDirection method of ColumnOrder.
func Test_ColumnOrder_GetDirection(t *testing.T) {
	o := ColumnOrder{
		Direction: contract.SortDesc,
	}
	assert.Equal(t, contract.SortDesc, o.GetDirection())
}

// Test_ColumnOrder_GetNulls tests the GetNulls method of ColumnOrder.
func Test_ColumnOrder_GetNulls(t *testing.T) {
	o := ColumnOrder{
		Nulls: contract.NullsLast,
	}
	assert.Equal(t, contract.NullsLast, o.GetNulls())
}

// Test_ColumnOrder_GetUsingOperator tests the GetUsingOperator method of
// ColumnOrder.
func Test_ColumnOrder_GetUsingOperator(t *testing.T) {
	o := ColumnOrder{
		UsingOperator: contract.OpEqual,
	}
	assert.Equal(t, contract.OpEqual, o.GetUsingOperator())
}
