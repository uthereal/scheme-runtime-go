package column

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func Test_IntervalColumn_Gt(t *testing.T) {
	c := IntervalColumn[any, pgtype.Interval]{
		Column: Column[any, pgtype.Interval]{Name: "duration"},
	}
	w := c.Gt(pgtype.Interval{Days: 1, Valid: true})
	assert.NotNil(t, w)
}

func Test_IntervalColumn_Gte(t *testing.T) {
	c := IntervalColumn[any, pgtype.Interval]{
		Column: Column[any, pgtype.Interval]{Name: "duration"},
	}
	w := c.Gte(pgtype.Interval{Days: 1, Valid: true})
	assert.NotNil(t, w)
}

func Test_IntervalColumn_Lt(t *testing.T) {
	c := IntervalColumn[any, pgtype.Interval]{
		Column: Column[any, pgtype.Interval]{Name: "duration"},
	}
	w := c.Lt(pgtype.Interval{Days: 1, Valid: true})
	assert.NotNil(t, w)
}

func Test_IntervalColumn_Lte(t *testing.T) {
	c := IntervalColumn[any, pgtype.Interval]{
		Column: Column[any, pgtype.Interval]{Name: "duration"},
	}
	w := c.Lte(pgtype.Interval{Days: 1, Valid: true})
	assert.NotNil(t, w)
}

func Test_NullableIntervalColumn_IsNull(t *testing.T) {
	c := NullableIntervalColumn[any, pgtype.Interval]{
		IntervalColumn: IntervalColumn[any, pgtype.Interval]{
			Column: Column[any, pgtype.Interval]{Name: "duration"},
		},
	}
	assert.NotNil(t, c.IsNull())
	assert.NotNil(t, c.IsNotNull())
}

func Test_NullableIntervalColumn_EqPtr(t *testing.T) {
	c := NullableIntervalColumn[any, pgtype.Interval]{
		IntervalColumn: IntervalColumn[any, pgtype.Interval]{
			Column: Column[any, pgtype.Interval]{Name: "duration"},
		},
	}
	val := pgtype.Interval{Days: 2, Valid: true}
	assert.NotNil(t, c.EqPtr(&val))
	assert.NotNil(t, c.NeqPtr(&val))
	assert.NotNil(t, c.GtPtr(&val))
	assert.NotNil(t, c.GtePtr(&val))
	assert.NotNil(t, c.LtPtr(&val))
	assert.NotNil(t, c.LtePtr(&val))
}

func Test_NullableIntervalColumn_Ordering(t *testing.T) {
	c := NullableIntervalColumn[any, pgtype.Interval]{
		IntervalColumn: IntervalColumn[any, pgtype.Interval]{
			Column: Column[any, pgtype.Interval]{Name: "duration"},
		},
	}
	assert.NotNil(t, c.AscNullsFirst())
	assert.NotNil(t, c.AscNullsLast())
	assert.NotNil(t, c.DescNullsFirst())
	assert.NotNil(t, c.DescNullsLast())
}
