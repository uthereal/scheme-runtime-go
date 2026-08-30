package column

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_DurationColumn_Gt tests the Gt method of DurationColumn.
func Test_DurationColumn_Gt(t *testing.T) {
	c := DurationColumn[any, time.Duration]{
		Name: "dur",
	}
	w := c.Gt(time.Second)
	assert.NotNil(t, w)
}

// Test_DurationColumn_Gte tests the Gte method of DurationColumn.
func Test_DurationColumn_Gte(t *testing.T) {
	c := DurationColumn[any, time.Duration]{
		Name: "dur",
	}
	w := c.Gte(time.Second)
	assert.NotNil(t, w)
}

// Test_DurationColumn_Lt tests the Lt method of DurationColumn.
func Test_DurationColumn_Lt(t *testing.T) {
	c := DurationColumn[any, time.Duration]{
		Name: "dur",
	}
	w := c.Lt(time.Second)
	assert.NotNil(t, w)
}

// Test_DurationColumn_Lte tests the Lte method of DurationColumn.
func Test_DurationColumn_Lte(t *testing.T) {
	c := DurationColumn[any, time.Duration]{
		Name: "dur",
	}
	w := c.Lte(time.Second)
	assert.NotNil(t, w)
}
