package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_BindingsTracker_Bind tests the Bind method of bindingsTracker to ensure
// proper positional parameter tracking and state immutability.
func Test_BindingsTracker_Bind(t *testing.T) {
	t.Run("first bind", func(t *testing.T) {
		b := bindingsTracker{}
		placeholder, newTracker := b.Bind("first")

		assert.Equal(t, "$1", placeholder)
		assert.Equal(t, []any{"first"}, newTracker.values)
		assert.Empty(t, b.values)
	})

	t.Run("subsequent bind", func(t *testing.T) {
		b := bindingsTracker{values: []any{"first"}}
		placeholder, finalTracker := b.Bind("second")

		assert.Equal(t, "$2", placeholder)
		assert.Equal(t, []any{"first", "second"}, finalTracker.values)
		assert.Equal(t, []any{"first"}, b.values)
	})
}
