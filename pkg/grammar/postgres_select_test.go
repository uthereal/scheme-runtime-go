package grammar

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/order"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// Test_PostgresGrammar_CompileSelect_Basic tests compiling a basic SELECT
// query.
func Test_PostgresGrammar_CompileSelect_Basic(t *testing.T) {
	g := NewPostgresGrammar()
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		columns:    []string{"id", "email"},
	}

	sql, bindings := g.CompileSelect(state)
	expectedSql := `SELECT "id", "email" FROM "public"."users"`
	assert.Equal(t, expectedSql, sql)
	assert.Empty(t, bindings)
}

// Test_PostgresGrammar_CompileSelect_WithWheres tests compiling SELECT with
// WHERE clauses.
func Test_PostgresGrammar_CompileSelect_WithWheres(t *testing.T) {
	g := NewPostgresGrammar()
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		columns:    []string{"*"},
		wheres: []contract.Where{
			where.BasicWhere{
				Column:   "id",
				Operator: ">",
				Value:    int64(10),
				Boolean:  "AND",
			},
			where.NullWhere{
				Column:  "deleted_at",
				Not:     false,
				Boolean: "AND",
			},
			where.InWhere{
				Column:  "role",
				Values:  []any{"admin", "editor"},
				Not:     false,
				Boolean: "OR",
			},
		},
	}

	sql, bindings := g.CompileSelect(state)
	expectedSql := `SELECT * FROM "public"."users" WHERE "id" > $1 AND "deleted_at" IS NULL OR "role" IN ($2, $3)`
	assert.Equal(t, expectedSql, sql)
	assert.Len(t, bindings, 3)
}

// Test_PostgresGrammar_CompileSelect_WithNestedWhere tests compiling SELECT
// with nested WHERE clauses.
func Test_PostgresGrammar_CompileSelect_WithNestedWhere(t *testing.T) {
	g := NewPostgresGrammar()
	nestedState := mockQueryState{
		wheres: []contract.Where{
			where.BasicWhere{
				Column:   "age",
				Operator: ">=",
				Value:    18,
				Boolean:  "AND",
			},
			where.BasicWhere{
				Column:   "status",
				Operator: "=",
				Value:    "active",
				Boolean:  "AND",
			},
		},
	}
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		columns:    []string{"*"},
		wheres: []contract.Where{
			where.BasicWhere{
				Column:   "deleted_at",
				Operator: "=",
				Value:    nil,
				Boolean:  "AND",
			},
			where.NestedWhere{
				Query:   nestedState,
				Boolean: "AND",
			},
		},
	}

	sql, bindings := g.CompileSelect(state)
	expectedSql := `SELECT * FROM "public"."users" WHERE "deleted_at" = $1 AND ("age" >= $2 AND "status" = $3)`
	assert.Equal(t, expectedSql, sql)
	assert.Len(t, bindings, 3)
}

// Test_PostgresGrammar_CompileSelect_WithOrderingAndPagination tests compiling
// SELECT with ordering and pagination.
func Test_PostgresGrammar_CompileSelect_WithOrderingAndPagination(
	t *testing.T,
) {
	g := NewPostgresGrammar()
	limit := uint64(5)
	offset := uint64(10)
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		columns:    []string{"id"},
		orders: []contract.Order{
			order.ColumnOrder{
				Column:    "created_at",
				Direction: contract.SortDesc,
				Nulls:     contract.NullsLast,
			},
		},
		limit:  &limit,
		offset: &offset,
	}

	sql, _ := g.CompileSelect(state)
	expectedSql := `SELECT "id" FROM "public"."users" ORDER BY "created_at" DESC NULLS LAST LIMIT 5 OFFSET 10`
	assert.Equal(t, expectedSql, sql)
}

// Test_PostgresGrammar_CompileSelect_WithDistinctAndAggregate tests compiling
// aggregate queries.
func Test_PostgresGrammar_CompileSelect_WithDistinctAndAggregate(t *testing.T) {
	g := NewPostgresGrammar()
	aggState := &contract.AggregateState{
		Function: contract.AggregateFunction("count"),
		Column:   "id",
	}

	// 1. Regular Aggregate
	stateRegular := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		aggregate:  aggState,
	}
	sqlRegular, _ := g.CompileSelect(stateRegular)
	assert.Equal(t, `SELECT count("id") AS aggregate FROM "public"."users"`, sqlRegular)

	// 2. Distinct Aggregate
	stateDistinct := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		distinct:   true,
		aggregate:  aggState,
	}
	sqlDistinct, _ := g.CompileSelect(stateDistinct)
	assert.Equal(t, `SELECT count(DISTINCT "id") AS aggregate FROM "public"."users"`, sqlDistinct)

	// 3. Distinct with Columns
	stateCols := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		columns:    []string{"id", "email"},
		distinct:   true,
	}
	sqlCols, _ := g.CompileSelect(stateCols)
	assert.Equal(t, `SELECT DISTINCT "id", "email" FROM "public"."users"`, sqlCols)

	// 4. Distinct Star (empty columns)
	stateStar := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		columns:    nil,
		distinct:   true,
	}
	sqlStar, _ := g.CompileSelect(stateStar)
	assert.Equal(t, `SELECT DISTINCT * FROM "public"."users"`, sqlStar)
}

// Test_PostgresGrammar_CompileSelect_AllOrderTypes tests compiling SELECT with
// all order types.
func Test_PostgresGrammar_CompileSelect_AllOrderTypes(t *testing.T) {
	g := NewPostgresGrammar()
	subState := mockQueryState{
		schemaName: "public",
		tableName:  "roles",
		columns:    []string{"id"},
	}

	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		columns:    []string{"id"},
		orders: []contract.Order{
			order.ColumnOrder{
				Column:    "created_at",
				Direction: contract.SortDesc,
				Nulls:     contract.NullsLast,
			},
			order.SubQueryOrder{
				Query:     subState,
				Direction: contract.SortAsc,
				Nulls:     contract.NullsFirst,
			},
			order.RawOrder{
				Sql:       "random()",
				Direction: contract.SortAsc,
			},
		},
	}

	sql, _ := g.CompileSelect(state)
	expectedSql := `SELECT "id" FROM "public"."users" ORDER BY "created_at" DESC NULLS LAST, (SELECT "id" FROM "public"."roles") ASC NULLS FIRST, random() ASC`
	assert.Equal(t, expectedSql, sql)
}

// Test_PostgresGrammar_CompileSelect_AllWhereTypes tests compiling SELECT with
// all 11 Where types.
func Test_PostgresGrammar_CompileSelect_AllWhereTypes(t *testing.T) {
	g := NewPostgresGrammar()

	subQuery := mockQueryState{
		schemaName: "public",
		tableName:  "roles",
		columns:    []string{"id"},
	}

	nestedQuery := mockQueryState{
		wheres: []contract.Where{
			where.BasicWhere{
				Column:   "status",
				Operator: "=",
				Value:    "active",
				Boolean:  "AND",
			},
		},
	}

	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		columns:    []string{"*"},
		wheres: []contract.Where{
			// 1. BasicWhere
			where.BasicWhere{
				Column:   "id",
				Operator: ">",
				Value:    10,
				Boolean:  "AND",
			},
			// 2. NullWhere
			where.NullWhere{
				Column:  "deleted_at",
				Not:     false,
				Boolean: "AND",
			},
			// 3. InWhere
			where.InWhere{
				Column:  "role",
				Values:  []any{"admin"},
				Not:     false,
				Boolean: "AND",
			},
			// 4. BetweenWhere
			where.BetweenWhere{
				Column:  "age",
				Min:     18,
				Max:     30,
				Not:     false,
				Boolean: "AND",
			},
			// 5. NestedWhere
			where.NestedWhere{
				Query:   nestedQuery,
				Boolean: "AND",
			},
			// 6. RawWhere
			where.RawWhere{
				Sql:      "status = ?",
				Bindings: []any{"active"},
				Boolean:  "AND",
			},
			// 7. ColumnWhere
			where.ColumnWhere{
				First:    "updated_at",
				Operator: ">",
				Second:   "created_at",
				Boolean:  "AND",
			},
			// 8. DateWhere
			where.DateWhere{
				Column:   "created_at",
				Operator: "=",
				Value:    "2026-08-08",
				Type:     "date",
				Boolean:  "AND",
			},
			// 9. ExistsWhere
			where.ExistsWhere{
				Query:   subQuery,
				Not:     false,
				Boolean: "AND",
			},
			// 10. SubQueryWhere
			where.SubQueryWhere{
				Column:   "role_id",
				Operator: "=",
				Query:    subQuery,
				Boolean:  "AND",
			},
			// 11. JsonWhere
			where.JsonWhere{
				Column:   "metadata",
				Key:      "preferences.theme",
				Operator: "=",
				Value:    "dark",
				Boolean:  "AND",
			},
		},
	}

	sql, bindings := g.CompileSelect(state)
	expectedSql := `SELECT * FROM "public"."users" ` +
		`WHERE "id" > $1 ` +
		`AND "deleted_at" IS NULL ` +
		`AND "role" IN ($2) ` +
		`AND "age" BETWEEN $3 AND $4 ` +
		`AND ("status" = $5) ` +
		`AND status = $6 ` +
		`AND "updated_at" > "created_at" ` +
		`AND "created_at"::date = $7 ` +
		`AND EXISTS (SELECT "id" FROM "public"."roles") ` +
		`AND "role_id" = (SELECT "id" FROM "public"."roles") ` +
		`AND "metadata" ->> $8 = $9`

	assert.Equal(t, expectedSql, sql)
	assert.Len(t, bindings, 9)
}

// Test_PostgresGrammar_CompileSelect_WithGroups tests compiling select queries
// with GROUP BY columns.
func Test_PostgresGrammar_CompileSelect_WithGroups(t *testing.T) {
	g := NewPostgresGrammar()
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		columns:    []string{"role", "status"},
		groups:     []string{"role", "status"},
	}

	sql, _ := g.CompileSelect(state)
	expectedSql := `SELECT "role", "status" FROM "public"."users" GROUP BY "role", "status"`
	assert.Equal(t, expectedSql, sql)
}

// Test_PostgresGrammar_CompileSelect_WithSubQueryAndParams tests compiling
// SELECT with subquery.
func Test_PostgresGrammar_CompileSelect_WithSubQueryAndParams(t *testing.T) {
	g := NewPostgresGrammar()

	subQuery := mockQueryState{
		schemaName: "public",
		tableName:  "roles",
		columns:    []string{"id"},
		wheres: []contract.Where{
			where.BasicWhere{
				Column:   "name",
				Operator: "=",
				Value:    "admin",
				Boolean:  "AND",
			},
		},
	}

	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		columns:    []string{"*"},
		wheres: []contract.Where{
			where.BasicWhere{
				Column:   "status",
				Operator: "=",
				Value:    "active",
				Boolean:  "AND",
			},
			where.SubQueryWhere{
				Column:   "role_id",
				Operator: "=",
				Query:    subQuery,
				Boolean:  "AND",
			},
		},
	}

	sql, bindings := g.CompileSelect(state)
	expectedSql := `SELECT * FROM "public"."users" WHERE "status" = $1 AND "role_id" = (SELECT "id" FROM "public"."roles" WHERE "name" = $2)`
	assert.Equal(t, expectedSql, sql)
	assert.Len(t, bindings, 2)
}

// Test_PostgresGrammar_Compile_ConcurrencyAndSafety tests query compilation
// concurrency safety.
func Test_PostgresGrammar_Compile_ConcurrencyAndSafety(t *testing.T) {
	g := NewPostgresGrammar()
	const goroutines = 20
	const iterations = 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				stateSelect := mockQueryState{
					schemaName: "public",
					tableName:  "users",
					columns:    []string{"id"},
					wheres: []contract.Where{
						where.BasicWhere{
							Column:   "status",
							Operator: "=",
							Value:    "active",
							Boolean:  "AND",
						},
					},
				}
				_, bindingsSel := g.CompileSelect(stateSelect)
				assert.Len(t, bindingsSel, 1)
			}
		}()
	}

	wg.Wait()
}

// Test_PostgresGrammar_CompileSelect_WithNilWhere tests compiling SELECT query
// with nil where clauses.
func Test_PostgresGrammar_CompileSelect_WithNilWhere(t *testing.T) {
	g := NewPostgresGrammar()
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		columns:    []string{"*"},
		wheres: []contract.Where{
			nil,
			where.BasicWhere{
				Column:   "status",
				Operator: "=",
				Value:    "active",
				Boolean:  "AND",
			},
			nil,
		},
	}

	sql, bindings := g.CompileSelect(state)
	expectedSql := `SELECT * FROM "public"."users" WHERE "status" = $1`
	assert.Equal(t, expectedSql, sql)
	assert.Len(t, bindings, 1)
}

// Test_PostgresGrammar_CompileInsert_NonRectangular tests compiling
// non-rectangular INSERT.
func Test_PostgresGrammar_CompileInsert_NonRectangular(t *testing.T) {
	g := NewPostgresGrammar()
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
	}
	values := [][]contract.ColumnValue{
		{
			{Column: "email", Value: "alice@example.com"},
		},
		{
			{Column: "email", Value: "bob@example.com"},
			{Column: "role", Value: "admin"},
		},
	}

	sql, bindings := g.CompileInsert(state, values)
	expectedSql := `INSERT INTO "public"."users" ("email", "role") SELECT unnested.unnested_column_1, unnested.unnested_column_2 FROM UNNEST($1::text[], $2::text[]) AS unnested(unnested_column_1, unnested_column_2)`
	assert.Equal(t, expectedSql, sql)
	assert.Len(t, bindings, 2)
}

// Test_PostgresGrammar_Compile_CustomSortOperatorValidation tests compiling
// custom sort operators.
func Test_PostgresGrammar_Compile_CustomSortOperatorValidation(t *testing.T) {
	g := NewPostgresGrammar()

	orderClause := order.ColumnOrder{
		Column:        "id",
		Direction:     contract.SortUsing,
		UsingOperator: "<",
		Nulls:         contract.NullsFirst,
	}
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		columns:    []string{"id"},
		orders:     []contract.Order{orderClause},
	}

	sql, _ := g.CompileSelect(state)
	expectedSql := `SELECT "id" FROM "public"."users" ORDER BY "id" USING < NULLS FIRST`
	assert.Equal(t, expectedSql, sql)
}

// Test_PostgresGrammar_CompileSelect_JsonWhere tests compiling SELECT with JSON
// where clauses.
func Test_PostgresGrammar_CompileSelect_JsonWhere(t *testing.T) {
	g := NewPostgresGrammar()

	jsonCond := where.JsonWhere{
		Column:   "metadata",
		Key:      "preferences.theme",
		Operator: contract.OpEqual,
		Value:    "dark",
		Boolean:  contract.BoolAnd,
	}
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		columns:    []string{"*"},
		wheres:     []contract.Where{jsonCond},
	}

	sql, bindings := g.CompileSelect(state)
	expectedSql := `SELECT * FROM "public"."users" WHERE "metadata" ->> $1 = $2`
	assert.Equal(t, expectedSql, sql)
	assert.Len(t, bindings, 2)
}
