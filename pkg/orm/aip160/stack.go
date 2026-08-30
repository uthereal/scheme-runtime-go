package aip160

import (
	"errors"
	"fmt"
)

// Pop removes and returns the top-most item from the stack.
func Pop[R any](
	stack []R,
) (R, []R, error) {
	var zero R
	n := len(stack)
	if n == 0 {
		return zero, stack, errors.New("attempted to pop from an empty stack")
	}

	last := stack[n-1]
	updatedStack := stack[:n-1]
	return last, updatedStack, nil
}

// Combine merges a given number of conditions from the top of the stack
// using the provided merge function.
func Combine[R any](
	stack []R,
	count int,
	merge func(args ...R) (R, error),
) ([]R, error) {
	if count <= 1 {
		return stack, nil
	}

	parts := make([]R, count)
	currentStack := stack

	for i := count - 1; i >= 0; i-- {
		var item R
		var err error
		item, currentStack, err = Pop(currentStack)
		if err != nil {
			return nil, err
		}
		parts[i] = item
	}

	combined, err := merge(parts...)
	if err != nil {
		return nil, fmt.Errorf("failed to merge conditions -> %w", err)
	}

	return append(currentStack, combined), nil
}
