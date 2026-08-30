package where

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_NestedWhere_GetQuery tests GetQuery of NestedWhere.
func Test_NestedWhere_GetQuery(t *testing.T) {
	q := mockQueryStateProvider{}
	w := NestedWhere{
		Query: q,
	}
	assert.Equal(t, q, w.GetQuery())
}

// Test_NestedWhere_IsNegated tests IsNegated of NestedWhere.
func Test_NestedWhere_IsNegated(t *testing.T) {
	w := NestedWhere{
		Not: true,
	}
	assert.True(t, w.IsNegated())
}

// Test_NestedWhere_GetBoolean tests GetBoolean of NestedWhere.
func Test_NestedWhere_GetBoolean(t *testing.T) {
	w := NestedWhere{
		Boolean: contract.BoolAnd,
	}
	assert.Equal(t, contract.BoolAnd, w.GetBoolean())
}
