package column

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_StringColumn_Contains tests the Contains method of StringColumn.
func Test_StringColumn_Contains(t *testing.T) {
	c := StringColumn[any, string]{
		Name: "s",
	}
	w1 := c.Contains("foo").(contract.BasicWhere)
	assert.Equal(t, "%foo%", w1.GetValue())

	w2 := c.Contains(`foo_bar%baz\qux`).(contract.BasicWhere)
	assert.Equal(t, `%foo\_bar\%baz\\qux%`, w2.GetValue())
}

// Test_StringColumn_ContainsFold tests the ContainsFold method of StringColumn.
func Test_StringColumn_ContainsFold(t *testing.T) {
	c := StringColumn[any, string]{
		Name: "s",
	}
	w1 := c.ContainsFold("foo").(contract.BasicWhere)
	assert.Equal(t, "%foo%", w1.GetValue())

	w2 := c.ContainsFold(`foo_bar%baz\qux`).(contract.BasicWhere)
	assert.Equal(t, `%foo\_bar\%baz\\qux%`, w2.GetValue())
}

// Test_StringColumn_HasPrefix tests the HasPrefix method of StringColumn.
func Test_StringColumn_HasPrefix(t *testing.T) {
	c := StringColumn[any, string]{
		Name: "s",
	}
	w1 := c.HasPrefix("foo").(contract.BasicWhere)
	assert.Equal(t, "foo%", w1.GetValue())

	w2 := c.HasPrefix(`foo_bar%baz\qux`).(contract.BasicWhere)
	assert.Equal(t, `foo\_bar\%baz\\qux%`, w2.GetValue())
}

// Test_StringColumn_HasSuffix tests the HasSuffix method of StringColumn.
func Test_StringColumn_HasSuffix(t *testing.T) {
	c := StringColumn[any, string]{
		Name: "s",
	}
	w1 := c.HasSuffix("foo").(contract.BasicWhere)
	assert.Equal(t, "%foo", w1.GetValue())

	w2 := c.HasSuffix(`foo_bar%baz\qux`).(contract.BasicWhere)
	assert.Equal(t, `%foo\_bar\%baz\\qux`, w2.GetValue())
}

// Test_StringColumn_Like tests the Like method of StringColumn.
func Test_StringColumn_Like(t *testing.T) {
	c := StringColumn[any, string]{
		Name: "s",
	}
	w := c.Like(`raw_%pattern\`).(contract.BasicWhere)
	assert.Equal(t, `raw_%pattern\`, w.GetValue())
}

// Test_StringColumn_ILike tests the ILike method of StringColumn.
func Test_StringColumn_ILike(t *testing.T) {
	c := StringColumn[any, string]{
		Name: "s",
	}
	w := c.ILike(`raw_%pattern\`).(contract.BasicWhere)
	assert.Equal(t, `raw_%pattern\`, w.GetValue())
}

// Test_StringColumn_TextSearchMatch tests TextSearchMatch of StringColumn.
func Test_StringColumn_TextSearchMatch(t *testing.T) {
	c := StringColumn[any, string]{
		Name: "s",
	}
	assert.NotNil(t, c.TextSearchMatch("foo"))
}
