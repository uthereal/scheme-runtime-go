package column

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_NullableJSONColumn_IsNull tests the IsNull method of NullableJSONColumn.
func Test_NullableJSONColumn_IsNull(t *testing.T) {
	c := NullableJSONColumn[any, map[string]any]{
		Name: "j",
	}
	assert.NotNil(t, c.IsNull())
}

// Test_NullableJSONColumn_IsNotNull tests IsNotNull of NullableJSONColumn.
func Test_NullableJSONColumn_IsNotNull(t *testing.T) {
	c := NullableJSONColumn[any, map[string]any]{
		Name: "j",
	}
	assert.NotNil(t, c.IsNotNull())
}

// Test_NullableJSONColumn_EqPtr tests the EqPtr method of NullableJSONColumn.
func Test_NullableJSONColumn_EqPtr(t *testing.T) {
	c := NullableJSONColumn[any, map[string]any]{
		Name: "j",
	}
	val := map[string]any{"foo": "bar"}
	assert.NotNil(t, c.EqPtr(&val))
}

// Test_NullableJSONColumn_NeqPtr tests NeqPtr method of NullableJSONColumn.
func Test_NullableJSONColumn_NeqPtr(t *testing.T) {
	c := NullableJSONColumn[any, map[string]any]{
		Name: "j",
	}
	val := map[string]any{"foo": "bar"}
	assert.NotNil(t, c.NeqPtr(&val))
}

// Test_NullableJSONColumn_AscNullsFirst tests AscNullsFirst.
func Test_NullableJSONColumn_AscNullsFirst(t *testing.T) {
	c := NullableJSONColumn[any, map[string]any]{
		Name: "j",
	}
	assert.NotNil(t, c.AscNullsFirst())
}

// Test_NullableJSONColumn_AscNullsLast tests AscNullsLast.
func Test_NullableJSONColumn_AscNullsLast(t *testing.T) {
	c := NullableJSONColumn[any, map[string]any]{
		Name: "j",
	}
	assert.NotNil(t, c.AscNullsLast())
}

// Test_NullableJSONColumn_DescNullsFirst tests DescNullsFirst.
func Test_NullableJSONColumn_DescNullsFirst(t *testing.T) {
	c := NullableJSONColumn[any, map[string]any]{
		Name: "j",
	}
	assert.NotNil(t, c.DescNullsFirst())
}

// Test_NullableJSONColumn_DescNullsLast tests DescNullsLast.
func Test_NullableJSONColumn_DescNullsLast(t *testing.T) {
	c := NullableJSONColumn[any, map[string]any]{
		Name: "j",
	}
	assert.NotNil(t, c.DescNullsLast())
}
