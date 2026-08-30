package column

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

// Test_DecimalColumn_Gt tests the Gt method of DecimalColumn.
func Test_DecimalColumn_Gt(t *testing.T) {
	c := DecimalColumn[any, pgtype.Numeric]{
		Name: "dec",
	}
	w := c.Gt(pgtype.Numeric{})
	assert.NotNil(t, w)
}

// Test_DecimalColumn_Gte tests the Gte method of DecimalColumn.
func Test_DecimalColumn_Gte(t *testing.T) {
	c := DecimalColumn[any, pgtype.Numeric]{
		Name: "dec",
	}
	w := c.Gte(pgtype.Numeric{})
	assert.NotNil(t, w)
}

// Test_DecimalColumn_Lt tests the Lt method of DecimalColumn.
func Test_DecimalColumn_Lt(t *testing.T) {
	c := DecimalColumn[any, pgtype.Numeric]{
		Name: "dec",
	}
	w := c.Lt(pgtype.Numeric{})
	assert.NotNil(t, w)
}

// Test_DecimalColumn_Lte tests the Lte method of DecimalColumn.
func Test_DecimalColumn_Lte(t *testing.T) {
	c := DecimalColumn[any, pgtype.Numeric]{
		Name: "dec",
	}
	w := c.Lte(pgtype.Numeric{})
	assert.NotNil(t, w)
}
