package column

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_ArrayColumn_Contains tests the Contains method of ArrayColumn.
func Test_ArrayColumn_Contains(t *testing.T) {
	c := ArrayColumn[any, string]{
		Name: "arr",
	}
	w := c.Contains([]string{"foo"})
	assert.NotNil(t, w)
}

// Test_ArrayColumn_Overlaps tests the Overlaps method of ArrayColumn.
func Test_ArrayColumn_Overlaps(t *testing.T) {
	c := ArrayColumn[any, string]{
		Name: "arr",
	}
	w := c.Overlaps([]string{"foo"})
	assert.NotNil(t, w)
}

// Test_ArrayColumn_ContainedBy tests the ContainedBy method of ArrayColumn.
func Test_ArrayColumn_ContainedBy(t *testing.T) {
	c := ArrayColumn[any, string]{
		Name: "arr",
	}
	w := c.ContainedBy([]string{"foo"})
	assert.NotNil(t, w)
}

// Test_ArrayColumn_Concat tests the Concat method of ArrayColumn.
func Test_ArrayColumn_Concat(t *testing.T) {
	c := ArrayColumn[any, string]{
		Name: "arr",
	}
	w := c.Concat([]string{"foo"})
	assert.NotNil(t, w)
}
