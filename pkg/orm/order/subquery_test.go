package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_SubQueryOrder_GetQuery tests the GetQuery method of SubQueryOrder.
func Test_SubQueryOrder_GetQuery(t *testing.T) {
	type mockQuery struct {
		contract.QueryStateProvider
	}
	m := mockQuery{}
	o := SubQueryOrder{
		Query: m,
	}
	assert.Equal(t, m, o.GetQuery())
}

// Test_SubQueryOrder_GetDirection tests the GetDirection method of
// SubQueryOrder.
func Test_SubQueryOrder_GetDirection(t *testing.T) {
	o := SubQueryOrder{
		Direction: contract.SortAsc,
	}
	assert.Equal(t, contract.SortAsc, o.GetDirection())
}

// Test_SubQueryOrder_GetNulls tests the GetNulls method of SubQueryOrder.
func Test_SubQueryOrder_GetNulls(t *testing.T) {
	o := SubQueryOrder{
		Nulls: contract.NullsFirst,
	}
	assert.Equal(t, contract.NullsFirst, o.GetNulls())
}
