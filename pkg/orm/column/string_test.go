package column

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_StringColumn_Contains tests the Contains method of StringColumn.
func Test_StringColumn_Contains(t *testing.T) {
	c := StringColumn[any, string]{
		Name: "s",
	}
	assert.NotNil(t, c.Contains("foo"))
}

// Test_StringColumn_ContainsFold tests the ContainsFold method of StringColumn.
func Test_StringColumn_ContainsFold(t *testing.T) {
	c := StringColumn[any, string]{
		Name: "s",
	}
	assert.NotNil(t, c.ContainsFold("foo"))
}

// Test_StringColumn_HasPrefix tests the HasPrefix method of StringColumn.
func Test_StringColumn_HasPrefix(t *testing.T) {
	c := StringColumn[any, string]{
		Name: "s",
	}
	assert.NotNil(t, c.HasPrefix("foo"))
}

// Test_StringColumn_HasSuffix tests the HasSuffix method of StringColumn.
func Test_StringColumn_HasSuffix(t *testing.T) {
	c := StringColumn[any, string]{
		Name: "s",
	}
	assert.NotNil(t, c.HasSuffix("foo"))
}

// Test_StringColumn_Like tests the Like method of StringColumn.
func Test_StringColumn_Like(t *testing.T) {
	c := StringColumn[any, string]{
		Name: "s",
	}
	assert.NotNil(t, c.Like("foo"))
}

// Test_StringColumn_ILike tests the ILike method of StringColumn.
func Test_StringColumn_ILike(t *testing.T) {
	c := StringColumn[any, string]{
		Name: "s",
	}
	assert.NotNil(t, c.ILike("foo"))
}

// Test_StringColumn_TextSearchMatch tests TextSearchMatch of StringColumn.
func Test_StringColumn_TextSearchMatch(t *testing.T) {
	c := StringColumn[any, string]{
		Name: "s",
	}
	assert.NotNil(t, c.TextSearchMatch("foo"))
}
