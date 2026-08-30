package column

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_NullableColumn_IsNull tests the IsNull method of NullableColumn.
func Test_NullableColumn_IsNull(t *testing.T) {
	c := NullableColumn[any, string]{
		Name: "c",
	}
	assert.NotNil(t, c.IsNull())
}

// Test_NullableColumn_IsNotNull tests the IsNotNull method of NullableColumn.
func Test_NullableColumn_IsNotNull(t *testing.T) {
	c := NullableColumn[any, string]{
		Name: "c",
	}
	assert.NotNil(t, c.IsNotNull())
}

// Test_NullableColumn_EqPtr tests the EqPtr method of NullableColumn.
func Test_NullableColumn_EqPtr(t *testing.T) {
	c := NullableColumn[any, string]{
		Name: "c",
	}
	val := "foo"
	assert.NotNil(t, c.EqPtr(&val))
}

// Test_NullableColumn_NeqPtr tests the NeqPtr method of NullableColumn.
func Test_NullableColumn_NeqPtr(t *testing.T) {
	c := NullableColumn[any, string]{
		Name: "c",
	}
	val := "foo"
	assert.NotNil(t, c.NeqPtr(&val))
}

// Test_NullableColumn_AscNullsFirst tests AscNullsFirst of NullableColumn.
func Test_NullableColumn_AscNullsFirst(t *testing.T) {
	c := NullableColumn[any, string]{
		Name: "c",
	}
	assert.NotNil(t, c.AscNullsFirst())
}

// Test_NullableColumn_AscNullsLast tests AscNullsLast of NullableColumn.
func Test_NullableColumn_AscNullsLast(t *testing.T) {
	c := NullableColumn[any, string]{
		Name: "c",
	}
	assert.NotNil(t, c.AscNullsLast())
}

// Test_NullableColumn_DescNullsFirst tests DescNullsFirst of NullableColumn.
func Test_NullableColumn_DescNullsFirst(t *testing.T) {
	c := NullableColumn[any, string]{
		Name: "c",
	}
	assert.NotNil(t, c.DescNullsFirst())
}

// Test_NullableColumn_DescNullsLast tests DescNullsLast of NullableColumn.
func Test_NullableColumn_DescNullsLast(t *testing.T) {
	c := NullableColumn[any, string]{
		Name: "c",
	}
	assert.NotNil(t, c.DescNullsLast())
}
