package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// Test_PostgresGrammar_CompileInsert tests the compilation of insert
// statements.
func Test_PostgresGrammar_CompileInsert(t *testing.T) {
	g := NewPostgresGrammar()
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
	}

	// 1. Basic success with multiple rows
	values := [][]contract.ColumnValue{
		{
			{Column: "email", Value: "alice@example.com"},
			{Column: "role", Value: "user"},
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

	// 2. Empty values slice
	sqlEmpty, bindingsEmpty := g.CompileInsert(state, nil)
	assert.Empty(t, sqlEmpty)
	assert.Nil(t, bindingsEmpty)
}

// Test_PostgresGrammar_CompileInsertReturning tests the compilation of insert
// statements with returning columns.
func Test_PostgresGrammar_CompileInsertReturning(t *testing.T) {
	g := NewPostgresGrammar()
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
	}

	values := [][]contract.ColumnValue{
		{
			{Column: "email", Value: "alice@example.com"},
		},
	}

	// 1. Success with returning columns
	sql, bindings := g.CompileInsertReturning(
		state,
		values,
		[]string{"id", "email"},
	)
	expectedSql := `INSERT INTO "public"."users" ("email") SELECT unnested.unnested_column_1 FROM UNNEST($1::text[]) AS unnested(unnested_column_1) RETURNING "id", "email"`
	assert.Equal(t, expectedSql, sql)
	assert.Len(t, bindings, 1)

	// 2. Success with empty returning columns
	sqlNoReturn, bindingsNoReturn := g.CompileInsertReturning(
		state,
		values,
		nil,
	)
	expectedNoReturnSql := `INSERT INTO "public"."users" ("email") SELECT unnested.unnested_column_1 FROM UNNEST($1::text[]) AS unnested(unnested_column_1)`
	assert.Equal(t, expectedNoReturnSql, sqlNoReturn)
	assert.Len(t, bindingsNoReturn, 1)

	// 3. Empty values returning
	sqlEmpty, bindingsEmpty := g.CompileInsertReturning(
		state,
		nil,
		[]string{"id"},
	)
	assert.Empty(t, sqlEmpty)
	assert.Nil(t, bindingsEmpty)
}

// Test_PostgresGrammar_CompileUpdate tests the compilation of update
// statements.
func Test_PostgresGrammar_CompileUpdate(t *testing.T) {
	g := NewPostgresGrammar()
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
	}

	values := []contract.ColumnValue{
		{Column: "email", Value: "new@example.com"},
	}

	// 1. Basic success without wheres
	sql, bindings := g.CompileUpdate(state, values)
	expectedSql := `UPDATE "public"."users" SET "email" = $1`
	assert.Equal(t, expectedSql, sql)
	assert.Len(t, bindings, 1)

	// 2. Success with wheres
	stateWithWheres := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		wheres: []contract.Where{
			where.BasicWhere{
				Column:   "id",
				Operator: "=",
				Value:    int64(1),
				Boolean:  "AND",
			},
		},
	}
	sqlWheres, bindingsWheres := g.CompileUpdate(stateWithWheres, values)
	expectedWheresSql := `UPDATE "public"."users" SET "email" = $1 WHERE "id" = $2`
	assert.Equal(t, expectedWheresSql, sqlWheres)
	assert.Len(t, bindingsWheres, 2)

	// 3. Empty values slice
	sqlEmpty, bindingsEmpty := g.CompileUpdate(state, nil)
	assert.Empty(t, sqlEmpty)
	assert.Nil(t, bindingsEmpty)
}

// Test_PostgresGrammar_CompileUpdateReturning tests the compilation of update
// statements with returning columns.
func Test_PostgresGrammar_CompileUpdateReturning(t *testing.T) {
	g := NewPostgresGrammar()
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
	}

	values := []contract.ColumnValue{
		{Column: "email", Value: "new@example.com"},
	}

	// 1. Success with returning columns
	sql, bindings := g.CompileUpdateReturning(
		state,
		values,
		[]string{"id", "email"},
	)
	expectedSql := `UPDATE "public"."users" SET "email" = $1 RETURNING "id", "email"`
	assert.Equal(t, expectedSql, sql)
	assert.Len(t, bindings, 1)

	// 2. Success with empty returning columns
	sqlNoReturn, bindingsNoReturn := g.CompileUpdateReturning(
		state,
		values,
		nil,
	)
	expectedNoReturnSql := `UPDATE "public"."users" SET "email" = $1`
	assert.Equal(t, expectedNoReturnSql, sqlNoReturn)
	assert.Len(t, bindingsNoReturn, 1)

	// 3. Empty values returning
	sqlEmpty, bindingsEmpty := g.CompileUpdateReturning(
		state,
		nil,
		[]string{"id"},
	)
	assert.Empty(t, sqlEmpty)
	assert.Nil(t, bindingsEmpty)
}

// Test_PostgresGrammar_CompileDelete tests the compilation of delete
// statements.
func Test_PostgresGrammar_CompileDelete(t *testing.T) {
	g := NewPostgresGrammar()
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
	}

	// 1. Basic success without wheres
	sql, bindings := g.CompileDelete(state)
	expectedSql := `DELETE FROM "public"."users"`
	assert.Equal(t, expectedSql, sql)
	assert.Empty(t, bindings)

	// 2. Success with wheres
	stateWithWheres := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		wheres: []contract.Where{
			where.BasicWhere{
				Column:   "status",
				Operator: "=",
				Value:    "suspended",
				Boolean:  "AND",
			},
		},
	}
	sqlWheres, bindingsWheres := g.CompileDelete(stateWithWheres)
	expectedWheresSql := `DELETE FROM "public"."users" WHERE "status" = $1`
	assert.Equal(t, expectedWheresSql, sqlWheres)
	assert.Len(t, bindingsWheres, 1)
}

// Test_PostgresGrammar_CompileDeleteReturning tests the compilation of delete
// statements with returning columns.
func Test_PostgresGrammar_CompileDeleteReturning(t *testing.T) {
	g := NewPostgresGrammar()
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
	}

	// 1. Success with returning columns and no wheres
	sql, bindings := g.CompileDeleteReturning(state, []string{"id"})
	expectedSql := `DELETE FROM "public"."users" RETURNING "id"`
	assert.Equal(t, expectedSql, sql)
	assert.Empty(t, bindings)

	// 2. Success with returning columns and wheres
	stateWithWheres := mockQueryState{
		schemaName: "public",
		tableName:  "users",
		wheres: []contract.Where{
			where.BasicWhere{
				Column:   "status",
				Operator: "=",
				Value:    "suspended",
				Boolean:  "AND",
			},
		},
	}
	sqlWheres, bindingsWheres := g.CompileDeleteReturning(
		stateWithWheres,
		[]string{"id"},
	)
	expectedWheresSql := `DELETE FROM "public"."users" WHERE "status" = $1 RETURNING "id"`
	assert.Equal(t, expectedWheresSql, sqlWheres)
	assert.Len(t, bindingsWheres, 1)

	// 3. Success with empty returning columns
	sqlNoReturn, bindingsNoReturn := g.CompileDeleteReturning(state, nil)
	expectedNoReturnSql := `DELETE FROM "public"."users"`
	assert.Equal(t, expectedNoReturnSql, sqlNoReturn)
	assert.Empty(t, bindingsNoReturn)
}

// Test_PostgresGrammar_CompileUpsert tests the compilation of upsert
// statements.
func Test_PostgresGrammar_CompileUpsert(t *testing.T) {
	g := NewPostgresGrammar()
	state := mockQueryState{
		schemaName: "public",
		tableName:  "users",
	}

	values := [][]contract.ColumnValue{
		{
			{Column: "id", Value: int64(1)},
			{Column: "email", Value: "upsert@example.com"},
		},
	}

	// 1. Success with updates (DO UPDATE)
	sql, bindings := g.CompileUpsert(state, values, []string{"id"})
	expectedSql := `INSERT INTO "public"."users" ("id", "email") SELECT unnested.unnested_column_1, unnested.unnested_column_2 FROM UNNEST($1::bigint[], $2::text[]) AS unnested(unnested_column_1, unnested_column_2) ON CONFLICT ("id") DO UPDATE SET "email" = EXCLUDED."email"`
	assert.Equal(t, expectedSql, sql)
	assert.Len(t, bindings, 2)

	// 2. Success with no updates (DO NOTHING)
	sqlNothing, bindingsNothing := g.CompileUpsert(
		state,
		values,
		[]string{"id", "email"},
	)
	expectedNothingSql := `INSERT INTO "public"."users" ("id", "email") SELECT unnested.unnested_column_1, unnested.unnested_column_2 FROM UNNEST($1::bigint[], $2::text[]) AS unnested(unnested_column_1, unnested_column_2) ON CONFLICT ("id", "email") DO NOTHING`
	assert.Equal(t, expectedNothingSql, sqlNothing)
	assert.Len(t, bindingsNothing, 2)

	// 3. Empty values upsert
	sqlEmpty, bindingsEmpty := g.CompileUpsert(state, nil, []string{"id"})
	assert.Empty(t, sqlEmpty)
	assert.Nil(t, bindingsEmpty)
}
