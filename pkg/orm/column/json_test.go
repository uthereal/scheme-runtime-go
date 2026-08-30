package column

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_JSONColumn_Contains tests the Contains method of JSONColumn.
func Test_JSONColumn_Contains(t *testing.T) {
	c := JSONColumn[any, map[string]any]{
		Name: "j",
	}
	w := c.Contains("foo")
	assert.NotNil(t, w)
}

// Test_JSONColumn_HasKey tests the HasKey method of JSONColumn.
func Test_JSONColumn_HasKey(t *testing.T) {
	c := JSONColumn[any, map[string]any]{
		Name: "j",
	}
	w := c.HasKey("foo")
	assert.NotNil(t, w)
}

// Test_JSONColumn_KeyEq tests the KeyEq method of JSONColumn.
func Test_JSONColumn_KeyEq(t *testing.T) {
	c := JSONColumn[any, map[string]any]{
		Name: "j",
	}
	w := c.KeyEq("foo", "bar")
	assert.NotNil(t, w)
}
