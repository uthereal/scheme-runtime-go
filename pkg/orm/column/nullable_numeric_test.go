package column

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_NullableNumericColumn_IsNull tests the IsNull method of
// NullableNumericColumn.
func Test_NullableNumericColumn_IsNull(t *testing.T) {
	c := NullableNumericColumn[any, int]{
		Name: "n",
	}
	assert.NotNil(t, c.IsNull())
}

// Test_NullableNumericColumn_IsNotNull tests IsNotNull of
// NullableNumericColumn.
func Test_NullableNumericColumn_IsNotNull(t *testing.T) {
	c := NullableNumericColumn[any, int]{
		Name: "n",
	}
	assert.NotNil(t, c.IsNotNull())
}

// Test_NullableNumericColumn_EqPtr tests the EqPtr method of
// NullableNumericColumn.
func Test_NullableNumericColumn_EqPtr(t *testing.T) {
	c := NullableNumericColumn[any, int]{
		Name: "n",
	}
	val := 10
	assert.NotNil(t, c.EqPtr(&val))
}

// Test_NullableNumericColumn_NeqPtr tests NeqPtr method of
// NullableNumericColumn.
func Test_NullableNumericColumn_NeqPtr(t *testing.T) {
	c := NullableNumericColumn[any, int]{
		Name: "n",
	}
	val := 10
	assert.NotNil(t, c.NeqPtr(&val))
}

// Test_NullableNumericColumn_AscNullsFirst tests AscNullsFirst.
func Test_NullableNumericColumn_AscNullsFirst(t *testing.T) {
	c := NullableNumericColumn[any, int]{
		Name: "n",
	}
	assert.NotNil(t, c.AscNullsFirst())
}

// Test_NullableNumericColumn_AscNullsLast tests AscNullsLast.
func Test_NullableNumericColumn_AscNullsLast(t *testing.T) {
	c := NullableNumericColumn[any, int]{
		Name: "n",
	}
	assert.NotNil(t, c.AscNullsLast())
}

// Test_NullableNumericColumn_DescNullsFirst tests DescNullsFirst.
func Test_NullableNumericColumn_DescNullsFirst(t *testing.T) {
	c := NullableNumericColumn[any, int]{
		Name: "n",
	}
	assert.NotNil(t, c.DescNullsFirst())
}

// Test_NullableNumericColumn_DescNullsLast tests DescNullsLast.
func Test_NullableNumericColumn_DescNullsLast(t *testing.T) {
	c := NullableNumericColumn[any, int]{
		Name: "n",
	}
	assert.NotNil(t, c.DescNullsLast())
}

// Test_NullableNumericColumn_BindType tests the BindType method of
// NullableNumericColumn.
func Test_NullableNumericColumn_BindType(t *testing.T) {
	c := NullableNumericColumn[any, int]{
		Name: "n",
	}
	assert.Equal(t, 0, c.BindType())
}
