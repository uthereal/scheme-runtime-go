package column

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_NullableArrayColumn_IsNull tests the IsNull method of
// NullableArrayColumn.
func Test_NullableArrayColumn_IsNull(t *testing.T) {
	c := NullableArrayColumn[any, string]{
		Name: "a",
	}
	assert.NotNil(t, c.IsNull())
}

// Test_NullableArrayColumn_IsNotNull tests IsNotNull of NullableArrayColumn.
func Test_NullableArrayColumn_IsNotNull(t *testing.T) {
	c := NullableArrayColumn[any, string]{
		Name: "a",
	}
	assert.NotNil(t, c.IsNotNull())
}

// Test_NullableArrayColumn_EqPtr tests the EqPtr method of NullableArrayColumn.
func Test_NullableArrayColumn_EqPtr(t *testing.T) {
	c := NullableArrayColumn[any, string]{
		Name: "a",
	}
	val := []string{"foo"}
	assert.NotNil(t, c.EqPtr(&val))
}

// Test_NullableArrayColumn_NeqPtr tests NeqPtr method of NullableArrayColumn.
func Test_NullableArrayColumn_NeqPtr(t *testing.T) {
	c := NullableArrayColumn[any, string]{
		Name: "a",
	}
	val := []string{"foo"}
	assert.NotNil(t, c.NeqPtr(&val))
}

// Test_NullableArrayColumn_AscNullsFirst tests AscNullsFirst.
func Test_NullableArrayColumn_AscNullsFirst(t *testing.T) {
	c := NullableArrayColumn[any, string]{
		Name: "a",
	}
	assert.NotNil(t, c.AscNullsFirst())
}

// Test_NullableArrayColumn_AscNullsLast tests AscNullsLast.
func Test_NullableArrayColumn_AscNullsLast(t *testing.T) {
	c := NullableArrayColumn[any, string]{
		Name: "a",
	}
	assert.NotNil(t, c.AscNullsLast())
}

// Test_NullableArrayColumn_DescNullsFirst tests DescNullsFirst.
func Test_NullableArrayColumn_DescNullsFirst(t *testing.T) {
	c := NullableArrayColumn[any, string]{
		Name: "a",
	}
	assert.NotNil(t, c.DescNullsFirst())
}

// Test_NullableArrayColumn_DescNullsLast tests DescNullsLast.
func Test_NullableArrayColumn_DescNullsLast(t *testing.T) {
	c := NullableArrayColumn[any, string]{
		Name: "a",
	}
	assert.NotNil(t, c.DescNullsLast())
}
