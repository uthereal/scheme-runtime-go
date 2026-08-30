package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_RawOrder_GetSql tests the GetSql method of RawOrder.
func Test_RawOrder_GetSql(t *testing.T) {
	o := RawOrder{
		Sql: "random()",
	}
	assert.Equal(t, "random()", o.GetSql())
}

// Test_RawOrder_GetDirection tests the GetDirection method of RawOrder.
func Test_RawOrder_GetDirection(t *testing.T) {
	o := RawOrder{
		Direction: contract.SortAsc,
	}
	assert.Equal(t, contract.SortAsc, o.GetDirection())
}
