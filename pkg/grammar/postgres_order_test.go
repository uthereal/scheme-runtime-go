package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/order"
)

// Test_PostgresGrammar_CompileOrders_Empty tests compileOrders with empty
// orders.
func Test_PostgresGrammar_CompileOrders_Empty(t *testing.T) {
	g := NewPostgresGrammar()
	state := mockQueryState{}
	tracker := bindingsTracker{}

	sql, _ := g.compileOrders(state, tracker)
	assert.Empty(t, sql)
}

// Test_PostgresGrammar_CompileOrder_Panic tests compileOrder with an unknown
// type.
func Test_PostgresGrammar_CompileOrder_Panic(t *testing.T) {
	g := NewPostgresGrammar()
	tracker := bindingsTracker{}

	type unknownOrder struct{}
	assert.Panics(t, func() {
		_, _ = g.compileOrder(unknownOrder{}, tracker)
	})
}

// Test_PostgresGrammar_CompileRawOrder_EmptyDirection tests compileRawOrder
// with empty direction.
func Test_PostgresGrammar_CompileRawOrder_EmptyDirection(t *testing.T) {
	g := NewPostgresGrammar()
	o := order.RawOrder{
		Sql:       "id",
		Direction: contract.SortDirection(""),
	}
	sql := g.compileRawOrder(o)
	assert.Equal(t, "id", sql)
}

// Test_PostgresGrammar_CompileColumnOrder_EmptyDirectionAndNulls tests
// compileColumnOrder with empty direction and nulls.
func Test_PostgresGrammar_CompileColumnOrder_EmptyDirectionAndNulls(
	t *testing.T,
) {
	g := NewPostgresGrammar()
	o := order.ColumnOrder{
		Column:    "id",
		Direction: contract.SortDirection(""),
		Nulls:     contract.NullsOrder(""),
	}
	sql := g.compileColumnOrder(o)
	assert.Equal(t, `"id"`, sql)
}

// Test_PostgresGrammar_CompileSubQueryOrder tests compileSubQueryOrder.
func Test_PostgresGrammar_CompileSubQueryOrder(t *testing.T) {
	g := NewPostgresGrammar()
	tracker := bindingsTracker{}

	subState := mockQueryState{
		schemaName: "public",
		tableName:  "roles",
		columns:    []string{"id"},
	}

	o := order.SubQueryOrder{
		Query:     subState,
		Direction: contract.SortAsc,
		Nulls:     contract.NullsFirst,
	}

	sql, _ := g.compileSubQueryOrder(o, tracker)
	expectedSql := `(SELECT "id" FROM "public"."roles") ASC NULLS FIRST`
	assert.Equal(t, expectedSql, sql)
}
