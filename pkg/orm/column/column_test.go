package column

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_Column_ColumnName tests the ColumnName method of Column.
func Test_Column_ColumnName(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	assert.Equal(t, "id", c.ColumnName())
}

// Test_Column_Asc tests the Asc method of Column.
func Test_Column_Asc(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	o := c.Asc()
	assert.NotNil(t, o)
}

// Test_Column_Desc tests the Desc method of Column.
func Test_Column_Desc(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	o := c.Desc()
	assert.NotNil(t, o)
}

// Test_Column_Eq tests the Eq method of Column.
func Test_Column_Eq(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	w := c.Eq("foo")
	assert.NotNil(t, w)
}

// Test_Column_Neq tests the Neq method of Column.
func Test_Column_Neq(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	w := c.Neq("foo")
	assert.NotNil(t, w)
}

// Test_Column_In tests the In method of Column.
func Test_Column_In(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	w := c.In("foo", "bar")
	assert.NotNil(t, w)
}

// Test_Column_NotIn tests the NotIn method of Column.
func Test_Column_NotIn(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	w := c.NotIn("foo", "bar")
	assert.NotNil(t, w)
}

// Test_Column_InQuery tests the InQuery method of Column.
func Test_Column_InQuery(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	type mockQueryStateProvider struct {
		contract.QueryStateProvider
	}
	q := mockQueryStateProvider{}
	w := c.InQuery(q)
	assert.NotNil(t, w)
}

// Test_Column_NotInQuery tests the NotInQuery method of Column.
func Test_Column_NotInQuery(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	type mockQueryStateProvider struct {
		contract.QueryStateProvider
	}
	q := mockQueryStateProvider{}
	w := c.NotInQuery(q)
	assert.NotNil(t, w)
}

// Test_Column_Between tests the Between method of Column.
func Test_Column_Between(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	w := c.Between("a", "z")
	assert.NotNil(t, w)
}

// Test_Column_NotBetween tests the NotBetween method of Column.
func Test_Column_NotBetween(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	w := c.NotBetween("a", "z")
	assert.NotNil(t, w)
}
