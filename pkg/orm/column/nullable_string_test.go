package column

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_NullableStringColumn_IsNull tests the IsNull method of
// NullableStringColumn.
func Test_NullableStringColumn_IsNull(t *testing.T) {
	c := NullableStringColumn[any, string]{
		Name: "s",
	}
	assert.NotNil(t, c.IsNull())
}

// Test_NullableStringColumn_IsNotNull tests IsNotNull of NullableStringColumn.
func Test_NullableStringColumn_IsNotNull(t *testing.T) {
	c := NullableStringColumn[any, string]{
		Name: "s",
	}
	assert.NotNil(t, c.IsNotNull())
}

// Test_NullableStringColumn_EqPtr tests the EqPtr method of
// NullableStringColumn.
func Test_NullableStringColumn_EqPtr(t *testing.T) {
	c := NullableStringColumn[any, string]{
		Name: "s",
	}
	val := "foo"
	assert.NotNil(t, c.EqPtr(&val))
}

// Test_NullableStringColumn_NeqPtr tests NeqPtr method of NullableStringColumn.
func Test_NullableStringColumn_NeqPtr(t *testing.T) {
	c := NullableStringColumn[any, string]{
		Name: "s",
	}
	val := "foo"
	assert.NotNil(t, c.NeqPtr(&val))
}

// Test_NullableStringColumn_AscNullsFirst tests AscNullsFirst.
func Test_NullableStringColumn_AscNullsFirst(t *testing.T) {
	c := NullableStringColumn[any, string]{
		Name: "s",
	}
	assert.NotNil(t, c.AscNullsFirst())
}

// Test_NullableStringColumn_AscNullsLast tests AscNullsLast.
func Test_NullableStringColumn_AscNullsLast(t *testing.T) {
	c := NullableStringColumn[any, string]{
		Name: "s",
	}
	assert.NotNil(t, c.AscNullsLast())
}

// Test_NullableStringColumn_DescNullsFirst tests DescNullsFirst.
func Test_NullableStringColumn_DescNullsFirst(t *testing.T) {
	c := NullableStringColumn[any, string]{
		Name: "s",
	}
	assert.NotNil(t, c.DescNullsFirst())
}

// Test_NullableStringColumn_DescNullsLast tests DescNullsLast.
func Test_NullableStringColumn_DescNullsLast(t *testing.T) {
	c := NullableStringColumn[any, string]{
		Name: "s",
	}
	assert.NotNil(t, c.DescNullsLast())
}
