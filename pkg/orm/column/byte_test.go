package column

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_ByteColumn_PostgresCast verifies that PostgresCast returns the
// correct bytea cast suffix.
func Test_ByteColumn_PostgresCast(t *testing.T) {
	c := ByteColumn[any]{
		Column: Column[any, []byte]{Name: "payload"},
	}
	assert.Equal(t, "::bytea[]", c.PostgresCast())
}

// Test_ByteColumn_ColumnName verifies ColumnName returns the correct column name.
func Test_ByteColumn_ColumnName(t *testing.T) {
	c := ByteColumn[any]{
		Column: Column[any, []byte]{Name: "data"},
	}
	assert.Equal(t, "data", c.ColumnName())
}

// Test_ByteColumn_QueryHelpers verifies basic query builder helpers on ByteColumn.
func Test_ByteColumn_QueryHelpers(t *testing.T) {
	c := ByteColumn[any]{
		Column: Column[any, []byte]{Name: "blob"},
	}

	assert.NotNil(t, c.Asc())
	assert.NotNil(t, c.Desc())
	assert.NotNil(t, c.Eq([]byte("foo")))
	assert.NotNil(t, c.Neq([]byte("bar")))
	assert.NotNil(t, c.In([]byte("foo"), []byte("bar")))
	assert.NotNil(t, c.NotIn([]byte("foo"), []byte("bar")))
	assert.NotNil(t, c.Between([]byte("a"), []byte("z")))
	assert.NotNil(t, c.NotBetween([]byte("a"), []byte("z")))
	assert.False(t, c.IsArray())
}

// Test_ByteColumn_ToTypedSlice verifies ToTypedSlice conversion.
func Test_ByteColumn_ToTypedSlice(t *testing.T) {
	c := ByteColumn[any]{
		Column: Column[any, []byte]{Name: "blob"},
	}

	// Slices of []byte
	b1 := []byte("hello")
	b2 := []byte("world")
	slice := []any{b1, b2}
	typed := c.ToTypedSlice(slice)
	assert.Equal(t, [][]byte{b1, b2}, typed)

	// Slices of *[]byte
	ptrSlice := []any{&b1, &b2, nil}
	typedPtr := c.ToTypedSlice(ptrSlice)
	assert.Equal(t, []*[]byte{&b1, &b2, nil}, typedPtr)
}

// Test_NullableByteColumn verifies all nullable condition and order helpers.
func Test_NullableByteColumn(t *testing.T) {
	c := NullableByteColumn[any]{
		ByteColumn: ByteColumn[any]{
			Column: Column[any, []byte]{Name: "blob"},
		},
	}

	assert.Equal(t, "::bytea[]", c.PostgresCast())
	assert.Equal(t, "blob", c.ColumnName())
	assert.NotNil(t, c.IsNull())
	assert.NotNil(t, c.IsNotNull())

	val := []byte("sample")
	assert.NotNil(t, c.EqPtr(&val))
	assert.NotNil(t, c.NeqPtr(&val))

	assert.NotNil(t, c.AscNullsFirst())
	assert.NotNil(t, c.AscNullsLast())
	assert.NotNil(t, c.DescNullsFirst())
	assert.NotNil(t, c.DescNullsLast())
}
