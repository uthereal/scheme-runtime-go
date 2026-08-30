package aip160

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.chromium.org/luci/common/data/aip160"
)

// Test_IterativeWalk_EmptyExpression tests IterativeWalk with an empty
// expression.
func Test_IterativeWalk_EmptyExpression(t *testing.T) {
	root := &aip160.Filter{}
	err := IterativeWalk(root, func(node any) error {
		return nil
	})
	assert.NoError(t, err)
}

// Test_IterativeWalk_Success tests iterative post-order traversal with
// callback.
func Test_IterativeWalk_Success(t *testing.T) {
	f, err := aip160.ParseFilter(
		"id = 10 OR (name = \"foo\" AND -is_active = true)",
	)
	assert.NoError(t, err)

	var visited []any
	err = IterativeWalk(f, func(node any) error {
		visited = append(visited, node)
		return nil
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, visited)
}

// Test_IterativeWalk_Error tests IterativeWalk when callback returns an error.
func Test_IterativeWalk_Error(t *testing.T) {
	f, err := aip160.ParseFilter("id = 10")
	assert.NoError(t, err)

	expectedErr := errors.New("walk error")
	err = IterativeWalk(f, func(node any) error {
		return expectedErr
	})
	assert.ErrorIs(t, err, expectedErr)
}
