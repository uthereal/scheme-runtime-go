package column

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_GeoColumn_Overlaps tests the Overlaps method of GeoColumn.
func Test_GeoColumn_Overlaps(t *testing.T) {
	c := GeoColumn[any, string]{
		Name: "g",
	}
	assert.NotNil(t, c.Overlaps("foo"))
}

// Test_GeoColumn_Contains tests the Contains method of GeoColumn.
func Test_GeoColumn_Contains(t *testing.T) {
	c := GeoColumn[any, string]{
		Name: "g",
	}
	assert.NotNil(t, c.Contains("foo"))
}

// Test_GeoColumn_ContainedBy tests the ContainedBy method of GeoColumn.
func Test_GeoColumn_ContainedBy(t *testing.T) {
	c := GeoColumn[any, string]{
		Name: "g",
	}
	assert.NotNil(t, c.ContainedBy("foo"))
}

// Test_GeoColumn_StrictLeft tests the StrictLeft method of GeoColumn.
func Test_GeoColumn_StrictLeft(t *testing.T) {
	c := GeoColumn[any, string]{
		Name: "g",
	}
	assert.NotNil(t, c.StrictLeft("foo"))
}

// Test_GeoColumn_StrictRight tests the StrictRight method of GeoColumn.
func Test_GeoColumn_StrictRight(t *testing.T) {
	c := GeoColumn[any, string]{
		Name: "g",
	}
	assert.NotNil(t, c.StrictRight("foo"))
}

// Test_GeoColumn_Below tests the Below method of GeoColumn.
func Test_GeoColumn_Below(t *testing.T) {
	c := GeoColumn[any, string]{
		Name: "g",
	}
	assert.NotNil(t, c.Below("foo"))
}

// Test_GeoColumn_Above tests the Above method of GeoColumn.
func Test_GeoColumn_Above(t *testing.T) {
	c := GeoColumn[any, string]{
		Name: "g",
	}
	assert.NotNil(t, c.Above("foo"))
}

// Test_GeoColumn_Distance tests the Distance method of GeoColumn.
func Test_GeoColumn_Distance(t *testing.T) {
	c := GeoColumn[any, string]{
		Name: "g",
	}
	assert.NotNil(t, c.Distance("foo"))
}

// Test_GeoColumn_ClosestProx tests the ClosestProx method of GeoColumn.
func Test_GeoColumn_ClosestProx(t *testing.T) {
	c := GeoColumn[any, string]{
		Name: "g",
	}
	assert.NotNil(t, c.ClosestProx("foo"))
}

// Test_NullableGeoColumn_IsNull tests the IsNull method of NullableGeoColumn.
func Test_NullableGeoColumn_IsNull(t *testing.T) {
	c := NullableGeoColumn[any, string]{
		Name: "g",
	}
	assert.NotNil(t, c.IsNull())
}

// Test_NullableGeoColumn_IsNotNull tests the IsNotNull method of
// NullableGeoColumn.
func Test_NullableGeoColumn_IsNotNull(t *testing.T) {
	c := NullableGeoColumn[any, string]{
		Name: "g",
	}
	assert.NotNil(t, c.IsNotNull())
}

// Test_NullableGeoColumn_EqPtr tests the EqPtr method of NullableGeoColumn.
func Test_NullableGeoColumn_EqPtr(t *testing.T) {
	c := NullableGeoColumn[any, string]{
		Name: "g",
	}
	val := "foo"
	assert.NotNil(t, c.EqPtr(&val))
}

// Test_NullableGeoColumn_NeqPtr tests the NeqPtr method of NullableGeoColumn.
func Test_NullableGeoColumn_NeqPtr(t *testing.T) {
	c := NullableGeoColumn[any, string]{
		Name: "g",
	}
	val := "foo"
	assert.NotNil(t, c.NeqPtr(&val))
}

// Test_NullableGeoColumn_AscNullsFirst tests AscNullsFirst of
// NullableGeoColumn.
func Test_NullableGeoColumn_AscNullsFirst(t *testing.T) {
	c := NullableGeoColumn[any, string]{
		Name: "g",
	}
	assert.NotNil(t, c.AscNullsFirst())
}

// Test_NullableGeoColumn_AscNullsLast tests AscNullsLast of NullableGeoColumn.
func Test_NullableGeoColumn_AscNullsLast(t *testing.T) {
	c := NullableGeoColumn[any, string]{
		Name: "g",
	}
	assert.NotNil(t, c.AscNullsLast())
}

// Test_NullableGeoColumn_DescNullsFirst tests DescNullsFirst of
// NullableGeoColumn.
func Test_NullableGeoColumn_DescNullsFirst(t *testing.T) {
	c := NullableGeoColumn[any, string]{
		Name: "g",
	}
	assert.NotNil(t, c.DescNullsFirst())
}

// Test_NullableGeoColumn_DescNullsLast tests DescNullsLast of
// NullableGeoColumn.
func Test_NullableGeoColumn_DescNullsLast(t *testing.T) {
	c := NullableGeoColumn[any, string]{
		Name: "g",
	}
	assert.NotNil(t, c.DescNullsLast())
}
