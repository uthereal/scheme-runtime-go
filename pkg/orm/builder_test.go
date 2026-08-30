package orm

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/grammar"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/aip160/fields"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/order"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
	freightv1 "go.einride.tech/aip/proto/gen/einride/example/freight/v1"
)

// dummyDecimalSubject is a dummy decimal subject for testing.
type dummyDecimalSubject struct{}

// mockRelation is a mock relation implementation for testing.
type mockRelation struct{}

// Eq performs an equal comparison.
func (d dummyDecimalSubject) Eq(
	val pgtype.Numeric,
) contract.Where {
	return where.BasicWhere{
		Column:   "price",
		Operator: contract.OpEqual,
		Value:    val,
	}
}

// Neq performs a not-equal comparison.
func (d dummyDecimalSubject) Neq(
	val pgtype.Numeric,
) contract.Where {
	return where.BasicWhere{
		Column:   "price",
		Operator: contract.OpNotEqual,
		Value:    val,
	}
}

// Gt performs a greater-than comparison.
func (d dummyDecimalSubject) Gt(
	val pgtype.Numeric,
) contract.Where {
	return where.BasicWhere{
		Column:   "price",
		Operator: contract.OpGreaterThan,
		Value:    val,
	}
}

// Gte performs a greater-than-or-equal comparison.
func (d dummyDecimalSubject) Gte(
	val pgtype.Numeric,
) contract.Where {
	return where.BasicWhere{
		Column:   "price",
		Operator: contract.OpGreaterThanOrEqual,
		Value:    val,
	}
}

// Lt performs a less-than comparison.
func (d dummyDecimalSubject) Lt(
	val pgtype.Numeric,
) contract.Where {
	return where.BasicWhere{
		Column:   "price",
		Operator: contract.OpLessThan,
		Value:    val,
	}
}

// Lte performs a less-than-or-equal comparison.
func (d dummyDecimalSubject) Lte(
	val pgtype.Numeric,
) contract.Where {
	return where.BasicWhere{
		Column:   "price",
		Operator: contract.OpLessThanOrEqual,
		Value:    val,
	}
}

// GetRelationQuery returns the query state provider for the relation.
func (r mockRelation) GetRelationQuery() contract.QueryStateProvider {
	return nil
}

// EagerLoad implements the EagerLoad method for mockRelation.
func (r mockRelation) EagerLoad(
	_ context.Context,
	_ contract.DB,
	_ []testModel,
) error {
	return nil
}

// Test_QueryBuilder_ColumnsAST tests standard and nullable columns.
func Test_QueryBuilder_ColumnsAST(t *testing.T) {
	// 1. Base Column Check
	eqCond := testSchema.Email.Eq("test@google.com").(where.BasicWhere)
	assert.Equal(t, "email", eqCond.Column)
	assert.Equal(t, contract.OpEqual, eqCond.Operator)
	assert.Equal(t, "test@google.com", eqCond.Value)

	inCond := testSchema.Email.In("a", "b").(where.InWhere)
	assert.Equal(t, "email", inCond.Column)
	require.Len(t, inCond.Values, 2)
	assert.Equal(t, "a", inCond.Values[0])
	assert.Equal(t, "b", inCond.Values[1])
	assert.False(t, inCond.Not)

	// 2. Nullable Pointer Checks (EqPtr / NeqPtr)
	someAge := 25
	ageCondVal := testSchema.Age.EqPtr(&someAge).(where.BasicWhere)
	assert.Equal(t, "age", ageCondVal.Column)
	assert.Equal(t, contract.OpEqual, ageCondVal.Operator)
	assert.Equal(t, 25, ageCondVal.Value)

	// 3. String wildcard methods
	containsCond := testSchema.Email.Contains("google").(where.BasicWhere)
	assert.Equal(t, contract.OpLike, containsCond.Operator)
	assert.Equal(t, "%google%", containsCond.Value)

	prefixCond := testSchema.Email.HasPrefix("admin").(where.BasicWhere)
	assert.Equal(t, contract.OpLike, prefixCond.Operator)
	assert.Equal(t, "admin%", prefixCond.Value)
}

// Test_QueryBuilder_TerminalExecutions tests execution of queries.
func Test_QueryBuilder_TerminalExecutions(t *testing.T) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	// 1. Test Get & Dynamic Hydration via Reflection
	mRows := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(1), "alice@google.com"},
			{int64(2), "bob@google.com"},
		},
	}
	db := &mockDb{queryRows: mRows}

	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)

	results, err := qb.Where(testSchema.ID.Gt(1)).Get(ctx)
	require.NoError(t, err)
	require.Len(t, results, 2)

	assert.Equal(t, int64(1), results[0].ID)
	assert.Equal(t, "alice@google.com", results[0].Email)
	assert.Equal(t, int64(2), results[1].ID)
	assert.Equal(t, "bob@google.com", results[1].Email)

	// 2. Test Count aggregate scan
	db.queryRowValue = &mockRow{value: int64(42)}
	count, err := qb.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint64(42), count)

	// 3. Test Exists Check
	exists, err := qb.Exists(ctx)
	require.NoError(t, err)
	assert.True(t, exists)

	// 4. Test Paginate checks
	mRowsPagination := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(10), "p1@google.com"},
			{int64(11), "p2@google.com"},
			{int64(12), "p3@google.com"},
		},
	}
	db.queryRows = mRowsPagination
	paginated, err := qb.Paginate(ctx, 2, 0)
	require.NoError(t, err)
	assert.True(t, paginated.HasMore)
	require.Len(t, paginated.Items, 2)
}

// Test_QueryBuilder_TerminalExecutions_ErrorPaths tests query execution
// error flows.
func Test_QueryBuilder_TerminalExecutions_ErrorPaths(t *testing.T) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	// 1. Test Exists error path when Count fails
	db := &mockDb{
		queryRowValue: &mockRow{err: errors.New("mock count error")},
	}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	exists, err := qb.Exists(ctx)
	assert.Error(t, err)
	assert.False(t, exists)
	assert.Contains(t, err.Error(), "mock count error")

	// 2. Test Get error path when rows.Err() fails
	mRows := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(1), "alice@google.com"},
		},
		rowsErr: errors.New("mock rows.Err error"),
	}
	dbErr := &mockDb{queryRows: mRows}
	qbErr := NewQueryBuilder[testModel, testMutator](
		dbErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, err = qbErr.Get(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mock rows.Err error")

	// 3. Test PaginateAip158 error path when cloned.Get fails
	dbGetErr := &mockDb{queryErr: errors.New("mock query error")}
	qbGetErr := NewQueryBuilder[testModel, testMutator](
		dbGetErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	req := &freightv1.ListShippersRequest{PageSize: 5}
	_, err = qbGetErr.PaginateAip158(ctx, req, 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mock query error")

	// 4. Test Get scan error
	mRowsScanErr := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{"not-an-int", "alice@google.com"},
		},
	}
	dbGetScanErr := &mockDb{queryRows: mRowsScanErr}
	qbGetScanErr := NewQueryBuilder[testModel, testMutator](
		dbGetScanErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, err = qbGetScanErr.Get(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed scanning database row")

	// 5. Test PaginateAip158 token parsing error
	qbToken := NewQueryBuilder[testModel, testMutator](
		&mockDb{},
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	reqTokenErr := &freightv1.ListShippersRequest{
		PageSize:  2,
		PageToken: "invalid-token",
	}
	_, err = qbToken.PaginateAip158(ctx, reqTokenErr, 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse page token")
}

// Test_QueryBuilder_PaginateAip158 tests AIP-158 pagination implementation.
func Test_QueryBuilder_PaginateAip158(t *testing.T) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	db := &mockDb{}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)

	// 1. First Page - mock DB has more records than PageSize
	mRows1 := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(10), "p1@google.com"},
			{int64(11), "p2@google.com"},
			{int64(12), "p3@google.com"},
		},
	}
	db.queryRows = mRows1

	req1 := &freightv1.ListShippersRequest{
		PageSize:  2,
		PageToken: "",
	}

	res1, err := qb.PaginateAip158(ctx, req1, 10)
	require.NoError(t, err)
	require.Len(t, res1.Items, 2)
	assert.Equal(t, int64(10), res1.Items[0].ID)
	assert.Equal(t, int64(11), res1.Items[1].ID)
	assert.NotEmpty(t, res1.NextPageToken)

	// 2. Second Page - using the returned NextPageToken
	mRows2 := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(12), "p3@google.com"},
			{int64(13), "p4@google.com"},
		},
	}
	db.queryRows = mRows2

	req2 := &freightv1.ListShippersRequest{
		PageSize:  2,
		PageToken: res1.NextPageToken,
	}

	res2, err := qb.PaginateAip158(ctx, req2, 10)
	require.NoError(t, err)
	require.Len(t, res2.Items, 2)
	assert.Equal(t, int64(12), res2.Items[0].ID)
	assert.Equal(t, int64(13), res2.Items[1].ID)
	assert.Empty(t, res2.NextPageToken)

	// 3. Invalid Page Token
	reqInvalid := &freightv1.ListShippersRequest{
		PageSize:  2,
		PageToken: "invalid-token",
	}
	_, errInvalid := qb.PaginateAip158(ctx, reqInvalid, 10)
	require.Error(t, errInvalid)

	// 4. DB Query Error
	dbErr := &mockDb{queryErr: errors.New("aip158 get failed")}
	qbErr := NewQueryBuilder[testModel, testMutator](
		dbErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	req3 := &freightv1.ListShippersRequest{
		PageSize: 2,
	}
	_, errErr := qbErr.PaginateAip158(ctx, req3, 10)
	require.Error(t, errErr)

	// 5. Unspecified Page Size (defaults to defaultPageSize)
	mRowsDefaultSize := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(10), "p1@google.com"},
		},
	}
	db.queryRows = mRowsDefaultSize
	reqDefaultSize := &freightv1.ListShippersRequest{
		PageSize:  0,
		PageToken: "",
	}
	resDefault, err := qb.PaginateAip158(ctx, reqDefaultSize, 10)
	require.NoError(t, err)
	require.Len(t, resDefault.Items, 1)
}

// Test_QueryBuilder_OrderByAip132 tests ordering of AIP-132 standards.
func Test_QueryBuilder_OrderByAip132(t *testing.T) {
	g := grammar.NewPostgresGrammar()

	mapColumnNameToOrderFunc := map[string]func(
		direction contract.SortDirection,
	) contract.Order{
		"id": func(dir contract.SortDirection) contract.Order {
			return order.ColumnOrder{Column: "id", Direction: dir}
		},
		"email": func(dir contract.SortDirection) contract.Order {
			return order.ColumnOrder{Column: "email", Direction: dir}
		},
		"settings.profile.age": func(
			dir contract.SortDirection,
		) contract.Order {
			return order.RawOrder{
				Sql:       "(settings -> 'profile' ->> 'age')::INTEGER",
				Direction: dir,
			}
		},
	}

	successTests := []struct {
		name        string
		orderBy     string
		expectedSql string
	}{
		{
			name:        "Empty string",
			orderBy:     "",
			expectedSql: `SELECT "id", "email" FROM "public"."users"`,
		},
		{
			name:        "Whitespace string",
			orderBy:     "   ",
			expectedSql: `SELECT "id", "email" FROM "public"."users"`,
		},
		{
			name:        "Single column default/ascending",
			orderBy:     "id",
			expectedSql: `SELECT "id", "email" FROM "public"."users" ORDER BY "id" ASC`,
		},
		{
			name:        "Single column explicitly descending",
			orderBy:     "email desc",
			expectedSql: `SELECT "id", "email" FROM "public"."users" ORDER BY "email" DESC`,
		},
		{
			name:        "Multiple columns mixed",
			orderBy:     "email desc, id",
			expectedSql: `SELECT "id", "email" FROM "public"."users" ORDER BY "email" DESC, "id" ASC`,
		},
		{
			name:        "Nested JSON field path with raw cast",
			orderBy:     "settings.profile.age desc",
			expectedSql: `SELECT "id", "email" FROM "public"."users" ORDER BY (settings -> 'profile' ->> 'age')::INTEGER DESC`,
		},
	}

	errorTests := []struct {
		name          string
		orderBy       string
		expectedError string
	}{
		{
			name:          "Unsupported column",
			orderBy:       "password",
			expectedError: "sorting by field \"password\" is not supported",
		},
		{
			name:          "Syntax error",
			orderBy:       "id,,email",
			expectedError: "invalid order_by syntax",
		},
	}

	for _, tc := range successTests {
		t.Run(tc.name, func(t *testing.T) {
			qb := NewQueryBuilder[testModel, testMutator](
				&mockDb{},
				g,
				testTable,
				testModelHydrate,
				testMutatorDehydrate,
			)

			res, err := qb.OrderByAip132(tc.orderBy, mapColumnNameToOrderFunc)
			require.NoError(t, err)
			assert.NotNil(t, res)

			sql, _ := g.CompileSelect(qb)
			assert.Equal(t, tc.expectedSql, sql)
		})
	}

	for _, tc := range errorTests {
		t.Run(tc.name, func(t *testing.T) {
			qb := NewQueryBuilder[testModel, testMutator](
				&mockDb{},
				g,
				testTable,
				testModelHydrate,
				testMutatorDehydrate,
			)

			res, err := qb.OrderByAip132(tc.orderBy, mapColumnNameToOrderFunc)
			require.Error(t, err)
			assert.Nil(t, res)
			assert.Contains(t, err.Error(), tc.expectedError)
		})
	}
}

// Test_QueryBuilder_WhereAip160 tests the AIP-160 filter transpilation.
func Test_QueryBuilder_WhereAip160(t *testing.T) {
	g := grammar.NewPostgresGrammar()

	mapColumnNameToAip160Field := map[string]contract.Aip160Field[contract.Where]{
		"id": fields.NewIntegerField[
			int64,
			contract.Where,
		](testSchema.ID),
		"email": fields.NewStringField[
			string,
			contract.Where,
		](testSchema.Email),
		"price": fields.NewDecimalField[
			pgtype.Numeric,
			contract.Where,
		](dummyDecimalSubject{}),
	}

	successTests := []struct {
		name        string
		filter      string
		expectedSql string
	}{
		{
			name:        "Empty string",
			filter:      "",
			expectedSql: `SELECT "id", "email" FROM "public"."users"`,
		},
		{
			name:        "Whitespace string",
			filter:      "   ",
			expectedSql: `SELECT "id", "email" FROM "public"."users"`,
		},
		{
			name:        "Simple equality numeric",
			filter:      "id = 10",
			expectedSql: `SELECT "id", "email" FROM "public"."users" WHERE "id" = $1`,
		},
		{
			name:        "Simple equality string",
			filter:      `email = "p1@google.com"`,
			expectedSql: `SELECT "id", "email" FROM "public"."users" WHERE "email" = $1`,
		},
		{
			name:        "Compound AND",
			filter:      `id = 10 AND email = "p1@google.com"`,
			expectedSql: `SELECT "id", "email" FROM "public"."users" WHERE ("email" = $1 AND "id" = $2)`,
		},
		{
			name:        "Compound OR",
			filter:      `id = 10 OR email = "p1@google.com"`,
			expectedSql: `SELECT "id", "email" FROM "public"."users" WHERE ("email" = $1 OR "id" = $2)`,
		},
		{
			name:        "Negation",
			filter:      `NOT id = 10`,
			expectedSql: `SELECT "id", "email" FROM "public"."users" WHERE NOT ("id" = $1)`,
		},
		{
			name:        "Complex nested logic",
			filter:      `(id = 10 OR email = "p1@google.com") AND id < 20`,
			expectedSql: `SELECT "id", "email" FROM "public"."users" WHERE ("id" < $1 AND ("email" = $2 OR "id" = $3))`,
		},
		{
			name:        "Exact decimal comparison Gt",
			filter:      "price > 99.99",
			expectedSql: `SELECT "id", "email" FROM "public"."users" WHERE "price" > $1`,
		},
		{
			name:        "Exact decimal comparison Eq",
			filter:      "price = 10.5",
			expectedSql: `SELECT "id", "email" FROM "public"."users" WHERE "price" = $1`,
		},
	}

	errorTests := []struct {
		name          string
		filter        string
		expectedError string
	}{
		{
			name:          "Unsupported column",
			filter:        "password = 123",
			expectedError: "filter field \"password\" not supported",
		},
		{
			name:          "Unsupported operator for type",
			filter:        "id : 10",
			expectedError: "does not support operator",
		},
		{
			name:          "Syntax error",
			filter:        "id == 10",
			expectedError: "failed to parse filter",
		},
		{
			name:          "Global search error",
			filter:        "abc",
			expectedError: "not supported",
		},
		{
			name:          "Missing argument error",
			filter:        "id",
			expectedError: "missing argument for filter field \"id\"",
		},
		{
			name:          "Coercion error",
			filter:        `id = "abc"`,
			expectedError: "failed to build condition",
		},
	}

	for _, tc := range successTests {
		t.Run(tc.name, func(t *testing.T) {
			qb := NewQueryBuilder[testModel, testMutator](
				&mockDb{},
				g,
				testTable,
				testModelHydrate,
				testMutatorDehydrate,
			)

			res, err := qb.WhereAip160(
				tc.filter,
				mapColumnNameToAip160Field,
			)
			require.NoError(t, err)
			assert.NotNil(t, res)

			sql, _ := g.CompileSelect(qb)
			assert.Equal(t, tc.expectedSql, sql)
		})
	}

	for _, tc := range errorTests {
		t.Run(tc.name, func(t *testing.T) {
			qb := NewQueryBuilder[testModel, testMutator](
				&mockDb{},
				g,
				testTable,
				testModelHydrate,
				testMutatorDehydrate,
			)

			res, err := qb.WhereAip160(
				tc.filter,
				mapColumnNameToAip160Field,
			)
			require.Error(t, err)
			assert.Nil(t, res)
			assert.Contains(t, err.Error(), tc.expectedError)
		})
	}
}

// Test_QueryBuilder_Getters tests all basic getters of QueryBuilder.
func Test_QueryBuilder_Getters(t *testing.T) {
	qb := NewQueryBuilder[testModel, testMutator](
		&mockDb{},
		grammar.NewPostgresGrammar(),
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)

	assert.Equal(t, "public", qb.GetSchemaName())
	assert.Equal(t, "users", qb.GetTableName())
	assert.Equal(t, []string{"id", "email"}, qb.GetSelectedColumns())
	assert.False(t, qb.IsDistinct())
	assert.Nil(t, qb.GetAggregate())
	assert.Empty(t, qb.GetWheres())
	assert.Empty(t, qb.GetOrders())
	assert.Empty(t, qb.GetGroups())
	assert.Empty(t, qb.GetHavings())

	lim, limOk := qb.GetLimit()
	assert.False(t, limOk)
	assert.Equal(t, uint64(0), lim)

	off, offOk := qb.GetOffset()
	assert.False(t, offOk)
	assert.Equal(t, uint64(0), off)

	qb.Limit(10)
	lim, limOk = qb.GetLimit()
	assert.True(t, limOk)
	assert.Equal(t, uint64(10), lim)

	qb.Offset(5)
	off, offOk = qb.GetOffset()
	assert.True(t, offOk)
	assert.Equal(t, uint64(5), off)
}

// Test_QueryBuilder_With tests the With method.
func Test_QueryBuilder_With(t *testing.T) {
	qb := NewQueryBuilder[testModel, testMutator](
		&mockDb{},
		grammar.NewPostgresGrammar(),
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	rel := mockRelation{}
	res := qb.With(rel)
	assert.NotNil(t, res)
	assert.Len(t, res.relations, 1)
}

// Test_QueryBuilder_Fluent tests fluent API methods.
func Test_QueryBuilder_Fluent(t *testing.T) {
	qb := NewQueryBuilder[testModel, testMutator](
		&mockDb{},
		grammar.NewPostgresGrammar(),
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)

	qb2 := qb.Distinct()
	assert.True(t, qb2.IsDistinct())

	qb3 := qb.GroupBy("email", "id")
	assert.Equal(t, []string{"email", "id"}, qb3.GetGroups())

	qb4 := qb.Having(testSchema.ID.Gt(10))
	assert.Len(t, qb4.GetHavings(), 1)
}

// Test_QueryBuilder_OrderBy tests the OrderBy method.
func Test_QueryBuilder_OrderBy(t *testing.T) {
	qb := NewQueryBuilder[testModel, testMutator](
		&mockDb{},
		grammar.NewPostgresGrammar(),
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	o := order.ColumnOrder{Column: "id", Direction: contract.SortAsc}
	qb.OrderBy(o)
	assert.Len(t, qb.GetOrders(), 1)
}

// Test_QueryBuilder_First tests the First method.
func Test_QueryBuilder_First(t *testing.T) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	// Case 1: Success
	mRows := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(42), "first@google.com"},
		},
	}
	db := &mockDb{queryRows: mRows}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	res, err := qb.First(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(42), res.ID)
	assert.Equal(t, "first@google.com", res.Email)

	// Case 2: No records found
	mRowsEmpty := &mockRows{
		cols:    []string{"id", "email"},
		records: [][]any{},
	}
	dbEmpty := &mockDb{queryRows: mRowsEmpty}
	qbEmpty := NewQueryBuilder[testModel, testMutator](
		dbEmpty,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, errEmpty := qbEmpty.First(ctx)
	require.Error(t, errEmpty)
	assert.Equal(t, "no matching record found", errEmpty.Error())

	// Case 3: DB error
	dbErr := &mockDb{queryErr: errors.New("db select error")}
	qbErr := NewQueryBuilder[testModel, testMutator](
		dbErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, errDb := qbErr.First(ctx)
	require.Error(t, errDb)
	assert.Contains(t, errDb.Error(), "db select error")
}

// Test_QueryBuilder_MutatorToColumnValues tests the mutatorToColumnValues.
func Test_QueryBuilder_MutatorToColumnValues(t *testing.T) {
	qb := NewQueryBuilder[testModel, testMutator](
		&mockDb{},
		grammar.NewPostgresGrammar(),
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	emailStr := "test@email.com"
	mut := testMutator{
		Email: &emailStr,
	}
	vals := qb.mutatorToColumnValues(&mut)
	require.Len(t, vals, 1)
	assert.Equal(t, "email", vals[0].Column)
	assert.Equal(t, "test@email.com", vals[0].Value)
}

// Test_QueryBuilder_Insert tests Insert and InsertMany.
func Test_QueryBuilder_Insert(t *testing.T) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	// Case 1: Insert Success
	db := &mockDb{}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	idVal := int64(10)
	emailVal := "test@test.com"
	mut := testMutator{
		ID:    &idVal,
		Email: &emailVal,
	}
	err := qb.Insert(ctx, mut)
	require.NoError(t, err)

	// Case 2: Insert Error
	dbErr := &mockDb{execErr: errors.New("exec fail")}
	qbErr := NewQueryBuilder[testModel, testMutator](
		dbErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	errErr := qbErr.Insert(ctx, mut)
	require.Error(t, errErr)
	assert.Contains(t, errErr.Error(), "exec fail")
}

// Test_QueryBuilder_Upsert tests Upsert and UpsertMany.
func Test_QueryBuilder_Upsert(t *testing.T) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	// Case 1: Upsert Success
	db := &mockDb{}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	idVal := int64(10)
	mut := testMutator{
		ID: &idVal,
	}
	err := qb.Upsert(ctx, mut, "id")
	require.NoError(t, err)

	// Case 2: Upsert Error
	dbErr := &mockDb{execErr: errors.New("exec upsert fail")}
	qbErr := NewQueryBuilder[testModel, testMutator](
		dbErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	errErr := qbErr.Upsert(ctx, mut, "id")
	require.Error(t, errErr)
	assert.Contains(t, errErr.Error(), "exec upsert fail")
}

// Test_QueryBuilder_UpdateAndDelete tests Update and Delete methods.
func Test_QueryBuilder_UpdateAndDelete(t *testing.T) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	// 1. Update Success
	db := &mockDb{}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	idVal := int64(20)
	mut := testMutator{
		ID: &idVal,
	}
	errUpdate := qb.Update(ctx, mut)
	require.NoError(t, errUpdate)

	// 2. Update Error
	dbErr := &mockDb{execErr: errors.New("update failed")}
	qbErr := NewQueryBuilder[testModel, testMutator](
		dbErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	errUpdateErr := qbErr.Update(ctx, mut)
	require.Error(t, errUpdateErr)
	assert.Contains(t, errUpdateErr.Error(), "update failed")

	// 3. Delete Success
	errDelete := qb.Delete(ctx)
	require.NoError(t, errDelete)

	// 4. Delete Error
	dbErrDel := &mockDb{execErr: errors.New("delete failed")}
	qbErrDel := NewQueryBuilder[testModel, testMutator](
		dbErrDel,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	errDeleteErr := qbErrDel.Delete(ctx)
	require.Error(t, errDeleteErr)
	assert.Contains(t, errDeleteErr.Error(), "delete failed")
}

// Test_QueryBuilder_InsertReturning tests InsertReturning methods.
func Test_QueryBuilder_InsertReturning(t *testing.T) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	// Case 1: Success
	mRows := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(10), "insert@test.com"},
		},
	}
	db := &mockDb{queryRows: mRows}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	idVal := int64(10)
	mut := testMutator{ID: &idVal}
	res, err := qb.InsertReturning(ctx, mut)
	require.NoError(t, err)
	assert.Equal(t, int64(10), res.ID)
	assert.Equal(t, "insert@test.com", res.Email)

	// Case 2: Empty return
	mRowsEmpty := &mockRows{
		cols:    []string{"id", "email"},
		records: [][]any{},
	}
	dbEmpty := &mockDb{queryRows: mRowsEmpty}
	qbEmpty := NewQueryBuilder[testModel, testMutator](
		dbEmpty,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, errEmpty := qbEmpty.InsertReturning(ctx, mut)
	require.Error(t, errEmpty)
	assert.Contains(
		t,
		errEmpty.Error(),
		"insert returning failed -> no record returned",
	)

	// Case 3: DB Query error
	dbErr := &mockDb{queryErr: errors.New("query error")}
	qbErr := NewQueryBuilder[testModel, testMutator](
		dbErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, errQuery := qbErr.InsertReturning(ctx, mut)
	require.Error(t, errQuery)
	assert.Contains(t, errQuery.Error(), "query error")

	// Case 4: Scan error
	mRowsScanErr := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{"not-an-int", "test@test.com"},
		},
	}
	dbScanErr := &mockDb{queryRows: mRowsScanErr}
	qbScanErr := NewQueryBuilder[testModel, testMutator](
		dbScanErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, errScan := qbScanErr.InsertReturning(ctx, mut)
	require.Error(t, errScan)
}

// Test_QueryBuilder_UpsertReturning tests UpsertReturning methods.
func Test_QueryBuilder_UpsertReturning(t *testing.T) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	// Case 1: Success
	mRows := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(10), "upsert@test.com"},
		},
	}
	db := &mockDb{queryRows: mRows}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	idVal := int64(10)
	mut := testMutator{ID: &idVal}
	res, err := qb.UpsertReturning(ctx, mut, "id")
	require.NoError(t, err)
	assert.Equal(t, int64(10), res.ID)
	assert.Equal(t, "upsert@test.com", res.Email)

	// Case 2: Empty return
	mRowsEmpty := &mockRows{
		cols:    []string{"id", "email"},
		records: [][]any{},
	}
	dbEmpty := &mockDb{queryRows: mRowsEmpty}
	qbEmpty := NewQueryBuilder[testModel, testMutator](
		dbEmpty,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, errEmpty := qbEmpty.UpsertReturning(ctx, mut, "id")
	require.Error(t, errEmpty)
	assert.Contains(
		t,
		errEmpty.Error(),
		"upsert returning failed -> no record returned",
	)

	// Case 3: DB Query error
	dbErr := &mockDb{queryErr: errors.New("upsert query error")}
	qbErr := NewQueryBuilder[testModel, testMutator](
		dbErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, errQuery := qbErr.UpsertReturning(ctx, mut, "id")
	require.Error(t, errQuery)
	assert.Contains(t, errQuery.Error(), "upsert query error")

	// Case 4: Scan error
	mRowsScanErr := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{"not-an-int", "test@test.com"},
		},
	}
	dbScanErr := &mockDb{queryRows: mRowsScanErr}
	qbScanErr := NewQueryBuilder[testModel, testMutator](
		dbScanErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, errScan := qbScanErr.UpsertReturning(ctx, mut, "id")
	require.Error(t, errScan)
}

// Test_QueryBuilder_UpdateReturning tests UpdateReturning method.
func Test_QueryBuilder_UpdateReturning(t *testing.T) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	// Case 1: Success
	mRows := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(10), "update@test.com"},
		},
	}
	db := &mockDb{queryRows: mRows}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	idVal := int64(10)
	mut := testMutator{ID: &idVal}
	res, err := qb.UpdateReturning(ctx, mut)
	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.Equal(t, int64(10), res[0].ID)
	assert.Equal(t, "update@test.com", res[0].Email)

	// Case 2: DB Query error
	dbErr := &mockDb{queryErr: errors.New("update query error")}
	qbErr := NewQueryBuilder[testModel, testMutator](
		dbErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, errQuery := qbErr.UpdateReturning(ctx, mut)
	require.Error(t, errQuery)
	assert.Contains(t, errQuery.Error(), "update query error")

	// Case 3: Scan error
	mRowsScanErr := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{"not-an-int", "test@test.com"},
		},
	}
	dbScanErr := &mockDb{queryRows: mRowsScanErr}
	qbScanErr := NewQueryBuilder[testModel, testMutator](
		dbScanErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, errScan := qbScanErr.UpdateReturning(ctx, mut)
	require.Error(t, errScan)
}

// Test_QueryBuilder_DeleteReturning tests DeleteReturning method.
func Test_QueryBuilder_DeleteReturning(t *testing.T) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	// Case 1: Success
	mRows := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(10), "delete@test.com"},
		},
	}
	db := &mockDb{queryRows: mRows}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	res, err := qb.DeleteReturning(ctx)
	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.Equal(t, int64(10), res[0].ID)
	assert.Equal(t, "delete@test.com", res[0].Email)

	// Case 2: DB Query error
	dbErr := &mockDb{queryErr: errors.New("delete query error")}
	qbErr := NewQueryBuilder[testModel, testMutator](
		dbErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, errQuery := qbErr.DeleteReturning(ctx)
	require.Error(t, errQuery)
	assert.Contains(t, errQuery.Error(), "delete query error")

	// Case 3: Scan error
	mRowsScanErr := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{"not-an-int", "test@test.com"},
		},
	}
	dbScanErr := &mockDb{queryRows: mRowsScanErr}
	qbScanErr := NewQueryBuilder[testModel, testMutator](
		dbScanErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, errScan := qbScanErr.DeleteReturning(ctx)
	require.Error(t, errScan)
}

// Test_QueryBuilder_PaginateWithCount tests PaginateWithCount method.
func Test_QueryBuilder_PaginateWithCount(t *testing.T) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	// Case 1: Success
	db := &mockDb{
		queryRowValue: &mockRow{value: 100},
		queryRows: &mockRows{
			cols: []string{"id", "email"},
			records: [][]any{
				{int64(10), "p1@google.com"},
			},
		},
	}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	res, err := qb.PaginateWithCount(ctx, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(100), res.TotalCount)
	require.Len(t, res.Items, 1)
	assert.Equal(t, int64(10), res.Items[0].ID)

	// Case 2: Count Query error
	dbErrCount := &mockDb{
		queryRowValue: &mockRow{err: errors.New("count query error")},
	}
	qbErrCount := NewQueryBuilder[testModel, testMutator](
		dbErrCount,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, errCount := qbErrCount.PaginateWithCount(ctx, 10, 0)
	require.Error(t, errCount)
	assert.Contains(t, errCount.Error(), "count query error")

	// Case 3: Get Items Query error
	dbErrGet := &mockDb{
		queryRowValue: &mockRow{value: 100},
		queryErr:      errors.New("get query error"),
	}
	qbErrGet := NewQueryBuilder[testModel, testMutator](
		dbErrGet,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, errGet := qbErrGet.PaginateWithCount(ctx, 10, 0)
	require.Error(t, errGet)
	assert.Contains(t, errGet.Error(), "get query error")
}

// Test_QueryBuilder_Aggregates tests Min, Max, Avg, and Sum methods.
func Test_QueryBuilder_Aggregates(t *testing.T) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	// 1. Min Success
	db := &mockDb{queryRowValue: &mockRow{value: int64(5)}}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	minVal, err := qb.Min(ctx, testSchema.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(5), minVal)

	// Min Error
	dbErr := &mockDb{queryRowValue: &mockRow{err: errors.New("min err")}}
	qbErr := NewQueryBuilder[testModel, testMutator](
		dbErr,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, errMin := qbErr.Min(ctx, testSchema.ID)
	require.Error(t, errMin)
	assert.Contains(t, errMin.Error(), "min err")

	// 2. Max Success
	dbMax := &mockDb{queryRowValue: &mockRow{value: int64(20)}}
	qbMax := NewQueryBuilder[testModel, testMutator](
		dbMax,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	maxVal, err := qbMax.Max(ctx, testSchema.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(20), maxVal)

	// Max Error
	_, errMax := qbErr.Max(ctx, testSchema.ID)
	require.Error(t, errMax)
	assert.Contains(t, errMax.Error(), "min err")

	// 3. Avg Success
	dbAvg := &mockDb{queryRowValue: &mockRow{value: int64(12)}}
	qbAvg := NewQueryBuilder[testModel, testMutator](
		dbAvg,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	avgVal, err := qbAvg.Avg(ctx, testSchema.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(12), avgVal)

	// Avg Error
	_, errAvg := qbErr.Avg(ctx, testSchema.ID)
	require.Error(t, errAvg)
	assert.Contains(t, errAvg.Error(), "min err")

	// 4. Sum Success
	dbSum := &mockDb{queryRowValue: &mockRow{value: int64(50)}}
	qbSum := NewQueryBuilder[testModel, testMutator](
		dbSum,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	sumVal, err := qbSum.Sum(ctx, testSchema.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(50), sumVal)

	// Sum Error
	_, errSum := qbErr.Sum(ctx, testSchema.ID)
	require.Error(t, errSum)
	assert.Contains(t, errSum.Error(), "min err")
}

// Test_QueryBuilder_Get_RowIterationError tests row iteration error in Get.
func Test_QueryBuilder_Get_RowIterationError(t *testing.T) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	mRows := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(1), "alice@google.com"},
		},
		rowsErr: errors.New("iteration failed"),
	}
	db := &mockDb{queryRows: mRows}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, err := qb.Get(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "iteration failed")
}

// Test_QueryBuilder_InsertReturningMany_RowIterationError tests row error.
func Test_QueryBuilder_InsertReturningMany_RowIterationError(
	t *testing.T,
) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	mRows := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(1), "alice@google.com"},
		},
		rowsErr: errors.New("insert iter failed"),
	}
	db := &mockDb{queryRows: mRows}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	muts := []testMutator{{}}
	_, err := qb.InsertReturningMany(ctx, muts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "insert iter failed")
}

// Test_QueryBuilder_UpsertReturningMany_RowIterationError tests row error.
func Test_QueryBuilder_UpsertReturningMany_RowIterationError(
	t *testing.T,
) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	mRows := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(1), "alice@google.com"},
		},
		rowsErr: errors.New("upsert iter failed"),
	}
	db := &mockDb{queryRows: mRows}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	muts := []testMutator{{}}
	_, err := qb.UpsertReturningMany(ctx, muts, "id")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "upsert iter failed")
}

// Test_QueryBuilder_UpdateReturning_RowIterationError tests row error
// in update.
func Test_QueryBuilder_UpdateReturning_RowIterationError(
	t *testing.T,
) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	mRows := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(1), "alice@google.com"},
		},
		rowsErr: errors.New("update iter failed"),
	}
	db := &mockDb{queryRows: mRows}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	mut := testMutator{}
	_, err := qb.UpdateReturning(ctx, mut)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "update iter failed")
}

// Test_QueryBuilder_DeleteReturning_RowIterationError tests row error
// in delete.
func Test_QueryBuilder_DeleteReturning_RowIterationError(
	t *testing.T,
) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	mRows := &mockRows{
		cols: []string{"id", "email"},
		records: [][]any{
			{int64(1), "alice@google.com"},
		},
		rowsErr: errors.New("delete iter failed"),
	}
	db := &mockDb{queryRows: mRows}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, err := qb.DeleteReturning(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "delete iter failed")
}

// Test_QueryBuilder_Paginate_Error tests error path of Paginate.
func Test_QueryBuilder_Paginate_Error(t *testing.T) {
	ctx := context.Background()
	compiler := grammar.NewPostgresGrammar()

	db := &mockDb{queryErr: errors.New("paginate get failed")}
	qb := NewQueryBuilder[testModel, testMutator](
		db,
		compiler,
		testTable,
		testModelHydrate,
		testMutatorDehydrate,
	)
	_, err := qb.Paginate(ctx, 10, 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "paginate get failed")
}

// Test_Where_WithBoolean tests WithBoolean method on all 11 concrete
// Where structs.
func Test_Where_WithBoolean(t *testing.T) {
	// 1. BasicWhere
	var w1 contract.Where = where.BasicWhere{Boolean: contract.BoolAnd}
	res1 := w1.WithBoolean(contract.BoolOr)
	assert.Equal(t, contract.BoolOr, res1.(where.BasicWhere).Boolean)

	// 2. BetweenWhere
	var w2 contract.Where = where.BetweenWhere{Boolean: contract.BoolAnd}
	res2 := w2.WithBoolean(contract.BoolOr)
	assert.Equal(t, contract.BoolOr, res2.(where.BetweenWhere).Boolean)

	// 3. ColumnWhere
	var w3 contract.Where = where.ColumnWhere{Boolean: contract.BoolAnd}
	res3 := w3.WithBoolean(contract.BoolOr)
	assert.Equal(t, contract.BoolOr, res3.(where.ColumnWhere).Boolean)

	// 4. DateWhere
	var w4 contract.Where = where.DateWhere{Boolean: contract.BoolAnd}
	res4 := w4.WithBoolean(contract.BoolOr)
	assert.Equal(t, contract.BoolOr, res4.(where.DateWhere).Boolean)

	// 5. ExistsWhere
	var w5 contract.Where = where.ExistsWhere{Boolean: contract.BoolAnd}
	res5 := w5.WithBoolean(contract.BoolOr)
	assert.Equal(t, contract.BoolOr, res5.(where.ExistsWhere).Boolean)

	// 6. InWhere
	var w6 contract.Where = where.InWhere{Boolean: contract.BoolAnd}
	res6 := w6.WithBoolean(contract.BoolOr)
	assert.Equal(t, contract.BoolOr, res6.(where.InWhere).Boolean)

	// 7. JsonWhere
	var w7 contract.Where = where.JsonWhere{Boolean: contract.BoolAnd}
	res7 := w7.WithBoolean(contract.BoolOr)
	assert.Equal(t, contract.BoolOr, res7.(where.JsonWhere).Boolean)

	// 8. NestedWhere
	var w8 contract.Where = where.NestedWhere{Boolean: contract.BoolAnd}
	res8 := w8.WithBoolean(contract.BoolOr)
	assert.Equal(t, contract.BoolOr, res8.(where.NestedWhere).Boolean)

	// 9. NullWhere
	var w9 contract.Where = where.NullWhere{Boolean: contract.BoolAnd}
	res9 := w9.WithBoolean(contract.BoolOr)
	assert.Equal(t, contract.BoolOr, res9.(where.NullWhere).Boolean)

	// 10. RawWhere
	var w10 contract.Where = where.RawWhere{Boolean: contract.BoolAnd}
	res10 := w10.WithBoolean(contract.BoolOr)
	assert.Equal(t, contract.BoolOr, res10.(where.RawWhere).Boolean)

	// 11. SubQueryWhere
	var w11 contract.Where = where.SubQueryWhere{Boolean: contract.BoolAnd}
	res11 := w11.WithBoolean(contract.BoolOr)
	assert.Equal(t, contract.BoolOr, res11.(where.SubQueryWhere).Boolean)
}
