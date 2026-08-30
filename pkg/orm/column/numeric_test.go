package column

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_NumericColumn_Gt tests the Gt method of NumericColumn.
func Test_NumericColumn_Gt(t *testing.T) {
	c := NumericColumn[any, int]{
		Name: "n",
	}
	assert.NotNil(t, c.Gt(10))
}

// Test_NumericColumn_Gte tests the Gte method of NumericColumn.
func Test_NumericColumn_Gte(t *testing.T) {
	c := NumericColumn[any, int]{
		Name: "n",
	}
	assert.NotNil(t, c.Gte(10))
}

// Test_NumericColumn_Lt tests the Lt method of NumericColumn.
func Test_NumericColumn_Lt(t *testing.T) {
	c := NumericColumn[any, int]{
		Name: "n",
	}
	assert.NotNil(t, c.Lt(10))
}

// Test_NumericColumn_Lte tests the Lte method of NumericColumn.
func Test_NumericColumn_Lte(t *testing.T) {
	c := NumericColumn[any, int]{
		Name: "n",
	}
	assert.NotNil(t, c.Lte(10))
}
