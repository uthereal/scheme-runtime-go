package aip160

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_Pop_Empty tests popping from an empty stack.
func Test_Pop_Empty(t *testing.T) {
	var stack []int
	val, updated, err := Pop(stack)
	assert.Zero(t, val)
	assert.Equal(t, stack, updated)
	assert.Error(t, err)
}

// Test_Pop_Success tests popping from a non-empty stack.
func Test_Pop_Success(t *testing.T) {
	stack := []int{10, 20}
	val, updated, err := Pop(stack)
	assert.Equal(t, 20, val)
	assert.Equal(t, []int{10}, updated)
	assert.NoError(t, err)
}

// Test_Combine_CountLessOrEqualOne tests Combine with count <= 1.
func Test_Combine_CountLessOrEqualOne(t *testing.T) {
	stack := []int{10}
	updated, err := Combine(stack, 1, nil)
	assert.NoError(t, err)
	assert.Equal(t, stack, updated)
}

// Test_Combine_Success tests a successful Combine call.
func Test_Combine_Success(t *testing.T) {
	stack := []int{10, 20}
	merge := func(args ...int) (int, error) {
		return args[0] + args[1], nil
	}
	updated, err := Combine(stack, 2, merge)
	assert.NoError(t, err)
	assert.Equal(t, []int{30}, updated)
}

// Test_Combine_PopError tests Combine when a pop fails due to insufficient
// items.
func Test_Combine_PopError(t *testing.T) {
	stack := []int{10}
	merge := func(args ...int) (int, error) {
		return args[0], nil
	}
	updated, err := Combine(stack, 2, merge)
	assert.Error(t, err)
	assert.Nil(t, updated)
}

// Test_Combine_MergeError tests Combine when the merge function returns an
// error.
func Test_Combine_MergeError(t *testing.T) {
	stack := []int{10, 20}
	expectedErr := errors.New("merge error")
	merge := func(args ...int) (int, error) {
		return 0, expectedErr
	}
	updated, err := Combine(stack, 2, merge)
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, updated)
}
