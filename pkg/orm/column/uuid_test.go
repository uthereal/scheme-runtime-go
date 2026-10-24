package column

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Test_UUIDColumn_PostgresCast verifies that PostgresCast returns the
// correct UUID cast suffix.
func Test_UUIDColumn_PostgresCast(t *testing.T) {
	c := UUIDColumn[any]{
		Column: Column[any, uuid.UUID]{Name: "id"},
	}
	assert.Equal(t, "::uuid[]", c.PostgresCast())
}

// Test_UUIDColumn_ToTypedSlice verifies that ToTypedSlice correctly handles
// uuid slices and pointer slices.
func Test_UUIDColumn_ToTypedSlice(t *testing.T) {
	c := UUIDColumn[any]{
		Column: Column[any, uuid.UUID]{Name: "id"},
	}

	u1 := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	u2 := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")

	// 1. Slice of uuid.UUID
	vals := []any{u1, u2}
	typed := c.ToTypedSlice(vals)
	assert.Equal(t, []uuid.UUID{u1, u2}, typed)

	// 2. Slice of pointer uuids
	ptrVals := []any{&u1, &u2, nil}
	typedPtr := c.ToTypedSlice(ptrVals)
	assert.Equal(t, []*uuid.UUID{&u1, &u2, nil}, typedPtr)

	// 3. Slice of strings parsed to uuid.UUID
	strVals := []any{
		"123e4567-e89b-12d3-a456-426614174000",
		"123e4567-e89b-12d3-a456-426614174001",
	}
	typedFromStr := c.ToTypedSlice(strVals)
	assert.Equal(t, []uuid.UUID{u1, u2}, typedFromStr)

	// 4. Empty slice
	empty := c.ToTypedSlice([]any{})
	assert.Equal(t, []uuid.UUID{}, empty)
}
