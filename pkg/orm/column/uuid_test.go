package column

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_UUIDColumn_PostgresCast verifies that PostgresCast returns the
// correct UUID cast suffix.
func Test_UUIDColumn_PostgresCast(t *testing.T) {
	c := UUIDColumn[any]{
		Name: "id",
	}
	assert.Equal(t, "::uuid[]", c.PostgresCast())
}

// Test_UUIDColumn_ToTypedSlice verifies that ToTypedSlice correctly handles
// string slices.
func Test_UUIDColumn_ToTypedSlice(t *testing.T) {
	c := UUIDColumn[any]{
		Name: "id",
	}

	// 1. Slice of strings
	vals := []any{"uuid1", "uuid2"}
	typed := c.ToTypedSlice(vals)
	assert.Equal(t, []string{"uuid1", "uuid2"}, typed)

	// 2. Slice of pointer strings (which converts to []*string because
	// UUIDColumn wraps Column[Model, string])
	u1, u2 := "uuid1", "uuid2"
	ptrVals := []any{&u1, &u2, nil}
	typedPtr := c.ToTypedSlice(ptrVals)
	assert.Equal(t, []*string{&u1, &u2, nil}, typedPtr)
}
