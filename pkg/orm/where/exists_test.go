package where

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_ExistsWhere_GetQuery tests the GetQuery method of ExistsWhere.
func Test_ExistsWhere_GetQuery(t *testing.T) {
	q := mockQueryStateProvider{}
	w := ExistsWhere{
		Query: q,
	}
	assert.Equal(t, q, w.GetQuery())
}

// Test_ExistsWhere_IsNot tests the IsNot method of ExistsWhere.
func Test_ExistsWhere_IsNot(t *testing.T) {
	w := ExistsWhere{
		Not: true,
	}
	assert.True(t, w.IsNot())
}

// Test_ExistsWhere_GetBoolean tests the GetBoolean method of ExistsWhere.
func Test_ExistsWhere_GetBoolean(t *testing.T) {
	w := ExistsWhere{
		Boolean: contract.BoolAnd,
	}
	assert.Equal(t, contract.BoolAnd, w.GetBoolean())
}
