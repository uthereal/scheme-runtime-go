package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// Test_PostgresGrammar_CompileHavings tests compileHavings.
func Test_PostgresGrammar_CompileHavings(t *testing.T) {
	g := NewPostgresGrammar()
	tracker := bindingsTracker{}

	// 1. Empty Havings
	stateEmpty := mockQueryState{}
	sqlEmpty, _ := g.compileHavings(stateEmpty, tracker)
	assert.Empty(t, sqlEmpty)

	// 2. Standard Havings
	stateHaving := mockQueryState{
		havings: []contract.Where{
			where.BasicWhere{
				Column:   "sum(amount)",
				Operator: ">",
				Value:    100,
				Boolean:  "AND",
			},
		},
	}
	sqlHaving, newTracker := g.compileHavings(stateHaving, tracker)
	assert.Equal(t, `HAVING "sum(amount)" > $1`, sqlHaving)
	assert.Len(t, newTracker.values, 1)
}

// Test_PostgresGrammar_CompileWheresRaw_Panic tests unknown Where types.
type unknownWhere struct{}

func (unknownWhere) GetBoolean() contract.BooleanOperator {
	return ""
}

func (unknownWhere) WithBoolean(_ contract.BooleanOperator) contract.Where {
	return unknownWhere{}
}

func Test_PostgresGrammar_CompileWheresRaw_Panic(t *testing.T) {
	g := NewPostgresGrammar()
	tracker := bindingsTracker{}

	assert.Panics(t, func() {
		_, _ = g.compileWheresRaw([]contract.Where{unknownWhere{}}, tracker)
	})
}

// Test_PostgresGrammar_CompileNullWhere_IsNotNull tests NullWhere with IsNot
// true.
func Test_PostgresGrammar_CompileNullWhere_IsNotNull(t *testing.T) {
	g := NewPostgresGrammar()
	tracker := bindingsTracker{}
	w := where.NullWhere{
		Column:  "deleted_at",
		Not:     true,
		Boolean: "AND",
	}

	sql, _ := g.compileNullWhere(w, tracker)
	assert.Equal(t, `"deleted_at" IS NOT NULL`, sql)
}

// Test_PostgresGrammar_CompileInWhere_EmptyAndNot tests InWhere with empty
// slice.
func Test_PostgresGrammar_CompileInWhere_EmptyAndNot(t *testing.T) {
	g := NewPostgresGrammar()
	tracker := bindingsTracker{}

	// 1. Empty IN (not negation)
	wIn := where.InWhere{
		Column:  "id",
		Values:  nil,
		Not:     false,
		Boolean: "AND",
	}
	sqlIn, _ := g.compileInWhere(wIn, tracker)
	assert.Equal(t, "1 = 0", sqlIn)

	// 2. Empty NOT IN (negation)
	wNotIn := where.InWhere{
		Column:  "id",
		Values:  nil,
		Not:     true,
		Boolean: "AND",
	}
	sqlNotIn, _ := g.compileInWhere(wNotIn, tracker)
	assert.Equal(t, "1 = 1", sqlNotIn)

	// 3. Non-empty NOT IN
	wNotInFull := where.InWhere{
		Column:  "role",
		Values:  []any{"admin"},
		Not:     true,
		Boolean: "AND",
	}
	sqlNotInFull, newTracker := g.compileInWhere(wNotInFull, tracker)
	assert.Equal(t, `"role" NOT IN ($1)`, sqlNotInFull)
	assert.Len(t, newTracker.values, 1)
}

// Test_PostgresGrammar_CompileBetweenWhere tests BetweenWhere compilation.
func Test_PostgresGrammar_CompileBetweenWhere(t *testing.T) {
	g := NewPostgresGrammar()
	tracker := bindingsTracker{}

	// 1. Regular Between
	wBetween := where.BetweenWhere{
		Column:  "age",
		Min:     18,
		Max:     30,
		Not:     false,
		Boolean: "AND",
	}
	sqlBetween, tracker2 := g.compileBetweenWhere(wBetween, tracker)
	assert.Equal(t, `"age" BETWEEN $1 AND $2`, sqlBetween)
	assert.Equal(t, []any{18, 30}, tracker2.values)

	// 2. Not Between
	wNotBetween := where.BetweenWhere{
		Column:  "age",
		Min:     18,
		Max:     30,
		Not:     true,
		Boolean: "AND",
	}
	sqlNotBetween, tracker3 := g.compileBetweenWhere(wNotBetween, tracker)
	assert.Equal(t, `"age" NOT BETWEEN $1 AND $2`, sqlNotBetween)
	assert.Equal(t, []any{18, 30}, tracker3.values)
}

// Test_PostgresGrammar_CompileNestedWhere_IsNegated tests NestedWhere negation.
func Test_PostgresGrammar_CompileNestedWhere_IsNegated(t *testing.T) {
	g := NewPostgresGrammar()
	tracker := bindingsTracker{}

	subQuery := mockQueryState{
		wheres: []contract.Where{
			where.BasicWhere{
				Column:   "status",
				Operator: "=",
				Value:    "active",
				Boolean:  "AND",
			},
		},
	}

	w := where.NestedWhere{
		Query:   subQuery,
		Not:     true,
		Boolean: "AND",
	}

	sql, _ := g.compileNestedWhere(w, tracker)
	assert.Equal(t, `NOT ("status" = $1)`, sql)
}

// Test_PostgresGrammar_CompileRawWhere tests RawWhere compilation.
func Test_PostgresGrammar_CompileRawWhere(t *testing.T) {
	g := NewPostgresGrammar()
	tracker := bindingsTracker{}

	w := where.RawWhere{
		Sql:      "id > ? AND name = ?",
		Bindings: []any{10, "alice"},
		Boolean:  "AND",
	}

	sql, newTracker := g.compileRawWhere(w, tracker)
	assert.Equal(t, "id > $1 AND name = $2", sql)
	assert.Equal(t, []any{10, "alice"}, newTracker.values)
}

// Test_PostgresGrammar_CompileColumnWhere tests ColumnWhere compilation.
func Test_PostgresGrammar_CompileColumnWhere(t *testing.T) {
	g := NewPostgresGrammar()
	tracker := bindingsTracker{}

	w := where.ColumnWhere{
		First:    "updated_at",
		Operator: ">",
		Second:   "created_at",
		Boolean:  "AND",
	}

	sql, _ := g.compileColumnWhere(w, tracker)
	assert.Equal(t, `"updated_at" > "created_at"`, sql)
}

// Test_PostgresGrammar_CompileDateWhere tests DateWhere extract type
// compilation.
func Test_PostgresGrammar_CompileDateWhere(t *testing.T) {
	g := NewPostgresGrammar()
	tracker := bindingsTracker{}

	// 1. Type "date"
	wDate := where.DateWhere{
		Column:   "created_at",
		Operator: "=",
		Value:    "2026-08-08",
		Type:     "date",
		Boolean:  "AND",
	}
	sql, _ := g.compileDateWhere(wDate, tracker)
	assert.Equal(t, `"created_at"::date = $1`, sql)

	// 2. Type "year"
	wYear := where.DateWhere{
		Column:   "created_at",
		Operator: "=",
		Value:    2026,
		Type:     "year",
		Boolean:  "AND",
	}
	sql, _ = g.compileDateWhere(wYear, tracker)
	assert.Equal(t, `EXTRACT(YEAR FROM "created_at") = $1`, sql)

	// 3. Type "month"
	wMonth := where.DateWhere{
		Column:   "created_at",
		Operator: "=",
		Value:    8,
		Type:     "month",
		Boolean:  "AND",
	}
	sql, _ = g.compileDateWhere(wMonth, tracker)
	assert.Equal(t, `EXTRACT(MONTH FROM "created_at") = $1`, sql)

	// 4. Type "day"
	wDay := where.DateWhere{
		Column:   "created_at",
		Operator: "=",
		Value:    8,
		Type:     "day",
		Boolean:  "AND",
	}
	sql, _ = g.compileDateWhere(wDay, tracker)
	assert.Equal(t, `EXTRACT(DAY FROM "created_at") = $1`, sql)

	// 5. Type "time"
	wTime := where.DateWhere{
		Column:   "created_at",
		Operator: "=",
		Value:    "12:34:56",
		Type:     "time",
		Boolean:  "AND",
	}
	sql, _ = g.compileDateWhere(wTime, tracker)
	assert.Equal(t, `"created_at"::time = $1`, sql)

	// 6. Type "unknown" (panic)
	wUnknown := where.DateWhere{
		Column:   "created_at",
		Operator: "=",
		Value:    "foo",
		Type:     "unknown",
		Boolean:  "AND",
	}
	assert.Panics(t, func() {
		_, _ = g.compileDateWhere(wUnknown, tracker)
	})
}

// Test_PostgresGrammar_CompileExistsWhere tests ExistsWhere compilation.
func Test_PostgresGrammar_CompileExistsWhere(t *testing.T) {
	g := NewPostgresGrammar()
	tracker := bindingsTracker{}

	subQuery := mockQueryState{
		schemaName: "public",
		tableName:  "roles",
		columns:    []string{"1"},
	}

	// 1. Regular EXISTS
	wExists := where.ExistsWhere{
		Query:   subQuery,
		Not:     false,
		Boolean: "AND",
	}
	sql, _ := g.compileExistsWhere(wExists, tracker)
	assert.Equal(t, `EXISTS (SELECT "1" FROM "public"."roles")`, sql)

	// 2. NOT EXISTS
	wNotExists := where.ExistsWhere{
		Query:   subQuery,
		Not:     true,
		Boolean: "AND",
	}
	sql, _ = g.compileExistsWhere(wNotExists, tracker)
	assert.Equal(t, `NOT EXISTS (SELECT "1" FROM "public"."roles")`, sql)
}
