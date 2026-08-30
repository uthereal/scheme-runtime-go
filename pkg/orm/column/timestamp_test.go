package column

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_TimestampColumn_Gt tests the Gt method of TimestampColumn.
func Test_TimestampColumn_Gt(t *testing.T) {
	c := TimestampColumn[any, time.Time]{
		Name: "t",
	}
	assert.NotNil(t, c.Gt(time.Now()))
}

// Test_TimestampColumn_Gte tests the Gte method of TimestampColumn.
func Test_TimestampColumn_Gte(t *testing.T) {
	c := TimestampColumn[any, time.Time]{
		Name: "t",
	}
	assert.NotNil(t, c.Gte(time.Now()))
}

// Test_TimestampColumn_Lt tests the Lt method of TimestampColumn.
func Test_TimestampColumn_Lt(t *testing.T) {
	c := TimestampColumn[any, time.Time]{
		Name: "t",
	}
	assert.NotNil(t, c.Lt(time.Now()))
}

// Test_TimestampColumn_Lte tests the Lte method of TimestampColumn.
func Test_TimestampColumn_Lte(t *testing.T) {
	c := TimestampColumn[any, time.Time]{
		Name: "t",
	}
	assert.NotNil(t, c.Lte(time.Now()))
}
