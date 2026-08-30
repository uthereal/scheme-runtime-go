package orm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	ormaip160 "github.com/uthereal/scheme-runtime-go/pkg/orm/aip160"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
	"go.chromium.org/luci/common/data/aip132"
	"go.chromium.org/luci/common/data/aip160"
	"go.einride.tech/aip/pagination"
)

// QueryBuilder compiles and executes SQL queries for a specific Model.
type QueryBuilder[Model any, Mutator any] struct {
	// db is the database connection interface.
	db contract.DB
	// compiler is the SQL compiler to compile the query state.
	compiler contract.SQLCompiler
	// schema is the metadata for the Model table.
	schema contract.TableMetadata[Model]
	// wheres is the list of filters to apply to the query.
	wheres []contract.Where
	// orders is the sorting criteria for the query.
	orders []contract.Order
	// relations is the list of relations to eager load.
	relations []contract.Relation[Model]
	// limit is the maximum number of rows to return.
	limit *uint64
	// offset is the number of rows to skip before returning.
	offset *uint64
	// distinct indicates whether to select unique records.
	distinct bool
	// groups is the columns to group results by.
	groups []string
	// havings is the filter conditions for grouped results.
	havings []contract.Where
	// aggregate is the select aggregate function state if any.
	aggregate *contract.AggregateState
	// columns is the list of specific columns to select.
	columns []string
	// hydrate is the function to populate model from columns.
	hydrate func(model *Model, columns []string) []any
	// dehydrate is the function to extract columns/values from mutator.
	dehydrate func(mutator *Mutator) ([]string, []any)
}

// NewQueryBuilder creates a new QueryBuilder instance.
func NewQueryBuilder[Model any, Mutator any](
	db contract.DB,
	compiler contract.SQLCompiler,
	schema contract.TableMetadata[Model],
	hydrate func(model *Model, columns []string) []any,
	dehydrate func(mutator *Mutator) ([]string, []any),
) *QueryBuilder[Model, Mutator] {
	return &QueryBuilder[Model, Mutator]{
		db:        db,
		compiler:  compiler,
		schema:    schema,
		hydrate:   hydrate,
		dehydrate: dehydrate,
	}
}

// Clone creates a shallow copy of the QueryBuilder state.
func (qb *QueryBuilder[Model, Mutator]) Clone() *QueryBuilder[
	Model,
	Mutator,
] {
	return &QueryBuilder[Model, Mutator]{
		db:        qb.db,
		compiler:  qb.compiler,
		schema:    qb.schema,
		wheres:    append([]contract.Where(nil), qb.wheres...),
		orders:    append([]contract.Order(nil), qb.orders...),
		relations: append([]contract.Relation[Model](nil), qb.relations...),
		limit:     qb.limit,
		offset:    qb.offset,
		distinct:  qb.distinct,
		groups:    append([]string(nil), qb.groups...),
		havings:   append([]contract.Where(nil), qb.havings...),
		aggregate: qb.aggregate,
		columns:   append([]string(nil), qb.columns...),
		hydrate:   qb.hydrate,
		dehydrate: qb.dehydrate,
	}
}

// SetDB updates the database connection of the QueryBuilder.
func (qb *QueryBuilder[Model, Mutator]) SetDB(db contract.DB) {
	qb.db = db
}

// GetSchemaName returns the schema name of the model.
func (qb *QueryBuilder[Model, Mutator]) GetSchemaName() string {
	return qb.schema.SchemaName
}

// GetTableName returns the table name of the model.
func (qb *QueryBuilder[Model, Mutator]) GetTableName() string {
	return qb.schema.TableName
}

// GetDefaultColumns returns the default columns to select from the table.
func (qb *QueryBuilder[Model, Mutator]) GetDefaultColumns() []string {
	cols := make([]string, len(qb.schema.DefaultColumns))
	for i, col := range qb.schema.DefaultColumns {
		cols[i] = col.ColumnName()
	}
	return cols
}

// GetColumnCastAndTypedSlice retrieves the SQL cast suffix, converts
// the untyped slice of values to its strongly-typed counterpart, and returns
// whether it is an array.
func (qb *QueryBuilder[Model, Mutator]) GetColumnCastAndTypedSlice(
	colName string,
	slice []any,
) (string, any, bool) {
	for _, col := range qb.schema.DefaultColumns {
		if col.ColumnName() != colName {
			continue
		}
		return col.PostgresCast(), col.ToTypedSlice(slice), col.IsArray()
	}
	return "::text[]", slice, false
}

// GetSelectedColumns returns the selected columns for the query.
func (qb *QueryBuilder[Model, Mutator]) GetSelectedColumns() []string {
	if len(qb.columns) > 0 {
		return qb.columns
	}
	return qb.GetDefaultColumns()
}

// Select sets the specific columns to select from the table.
func (qb *QueryBuilder[Model, Mutator]) Select(
	columns ...contract.Column[Model],
) *QueryBuilder[Model, Mutator] {
	qb.columns = make([]string, len(columns))
	for i, col := range columns {
		qb.columns[i] = col.ColumnName()
	}
	return qb
}

// IsDistinct returns whether the query selects unique records.
func (qb *QueryBuilder[Model, Mutator]) IsDistinct() bool {
	return qb.distinct
}

// GetAggregate returns the aggregate state of the query or nil if not set.
func (qb *QueryBuilder[Model, Mutator]) GetAggregate(
) *contract.AggregateState {
	return qb.aggregate
}

// GetWheres returns the slice of where clauses applied to the query.
func (qb *QueryBuilder[Model, Mutator]) GetWheres() []contract.Where {
	return qb.wheres
}

// GetOrders returns the slice of order clauses applied to the query.
func (qb *QueryBuilder[Model, Mutator]) GetOrders() []contract.Order {
	return qb.orders
}

// GetGroups returns the slice of group by columns.
func (qb *QueryBuilder[Model, Mutator]) GetGroups() []string {
	return qb.groups
}

// GetHavings returns the slice of having clauses applied to the query.
func (qb *QueryBuilder[Model, Mutator]) GetHavings() []contract.Where {
	return qb.havings
}

// GetLimit returns the limit and true if it is set.
func (qb *QueryBuilder[Model, Mutator]) GetLimit() (
	uint64,
	bool,
) {
	if qb.limit == nil {
		return 0, false
	}
	return *qb.limit, true
}

// GetOffset returns the offset and true if it is set.
func (qb *QueryBuilder[Model, Mutator]) GetOffset() (
	uint64,
	bool,
) {
	if qb.offset == nil {
		return 0, false
	}
	return *qb.offset, true
}

// Limit sets the maximum number of records to return.
func (qb *QueryBuilder[Model, Mutator]) Limit(
	limit uint64,
) *QueryBuilder[Model, Mutator] {
	qb.limit = &limit
	return qb
}

// Offset sets the number of records to skip before starting to return rows.
func (qb *QueryBuilder[Model, Mutator]) Offset(
	offset uint64,
) *QueryBuilder[Model, Mutator] {
	qb.offset = &offset
	return qb
}

// Where appends any valid Where clause.
func (qb *QueryBuilder[Model, Mutator]) Where(
	w contract.Where,
) *QueryBuilder[Model, Mutator] {
	qb.wheres = append(qb.wheres, w)
	return qb
}

// WhereAip160 parses and appends an AIP-160 filter string to the query.
func (qb *QueryBuilder[Model, Mutator]) WhereAip160(
	filter string,
	mapColumnNameToAip160Field map[string]contract.Aip160Field[contract.Where],
) (*QueryBuilder[Model, Mutator], error) {
	if strings.TrimSpace(filter) == "" {
		return qb, nil
	}

	f, err := aip160.ParseFilter(filter)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse filter %q -> %w",
			filter,
			err,
		)
	}

	var stack []contract.Where

	err = ormaip160.IterativeWalk(f, func(node any) error {
		switch n := node.(type) {
		case *aip160.Restriction:
			fieldName := n.Comparable.Member.Input()
			field, ok := mapColumnNameToAip160Field[fieldName]
			if !ok {
				return fmt.Errorf(
					"filter field %q not supported",
					fieldName,
				)
			}

			op := contract.Aip160Operator(n.Comparator)
			if !ormaip160.HasOperator(field.Type(), op) {
				return fmt.Errorf(
					"filter field %q with type %q does not support operator %q",
					fieldName,
					field.Type(),
					op,
				)
			}

			if n.Arg == nil {
				return fmt.Errorf(
					"missing argument for filter field %q",
					fieldName,
				)
			}

			cond, err := field.BuildCondition(op, n.Arg)
			if err != nil {
				return fmt.Errorf(
					"failed to build condition for field %q -> %w",
					fieldName,
					err,
				)
			}
			stack = append(stack, cond)

		case *aip160.Term:
			if n.Negated {
				var item contract.Where
				var popErr error
				item, stack, popErr = ormaip160.Pop(stack)
				if popErr != nil {
					return popErr
				}
				subQB := NewQueryBuilder[Model, Mutator](
					nil,
					nil,
					contract.TableMetadata[Model]{},
					nil,
					nil,
				)
				subQB.wheres = []contract.Where{item}
				stack = append(
					stack,
					where.NestedWhere{
						Query:   subQB,
						Not:     true,
						Boolean: contract.BoolAnd,
					},
				)
			}

		case *aip160.Factor:
			var combineErr error
			stack, combineErr = ormaip160.Combine(
				stack,
				len(n.Terms),
				func(
					parts ...contract.Where,
				) (contract.Where, error) {
					for i := 1; i < len(parts); i++ {
						parts[i] = parts[i].WithBoolean(
							contract.BoolOr,
						)
					}
					subQB := NewQueryBuilder[Model, Mutator](
						nil,
						nil,
						contract.TableMetadata[Model]{},
						nil,
						nil,
					)
					subQB.wheres = parts
					return where.NestedWhere{
						Query:   subQB,
						Boolean: contract.BoolAnd,
					}, nil
				},
			)
			if combineErr != nil {
				return combineErr
			}

		case *aip160.Sequence:
			var combineErr error
			stack, combineErr = ormaip160.Combine(
				stack,
				len(n.Factors),
				func(
					parts ...contract.Where,
				) (contract.Where, error) {
					for i := 1; i < len(parts); i++ {
						parts[i] = parts[i].WithBoolean(
							contract.BoolAnd,
						)
					}
					subQB := NewQueryBuilder[Model, Mutator](
						nil,
						nil,
						contract.TableMetadata[Model]{},
						nil,
						nil,
					)
					subQB.wheres = parts
					return where.NestedWhere{
						Query:   subQB,
						Boolean: contract.BoolAnd,
					}, nil
				},
			)
			if combineErr != nil {
				return combineErr
			}

		case *aip160.Expression:
			var combineErr error
			stack, combineErr = ormaip160.Combine(
				stack,
				len(n.Sequences),
				func(
					parts ...contract.Where,
				) (contract.Where, error) {
					for i := 1; i < len(parts); i++ {
						parts[i] = parts[i].WithBoolean(
							contract.BoolAnd,
						)
					}
					subQB := NewQueryBuilder[Model, Mutator](
						nil,
						nil,
						contract.TableMetadata[Model]{},
						nil,
						nil,
					)
					subQB.wheres = parts
					return where.NestedWhere{
						Query:   subQB,
						Boolean: contract.BoolAnd,
					}, nil
				},
			)
			if combineErr != nil {
				return combineErr
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	qb.wheres = append(qb.wheres, stack...)
	return qb, nil
}

// OrderBy appends any valid Order clause.
func (qb *QueryBuilder[Model, Mutator]) OrderBy(
	o contract.Order,
) *QueryBuilder[Model, Mutator] {
	qb.orders = append(qb.orders, o)
	return qb
}

// OrderByAip132 parses and appends an AIP-132 order string to the query.
func (qb *QueryBuilder[Model, Mutator]) OrderByAip132(
	orderBy string,
	mapColumnNameToOrderFunc map[string]func(
		direction contract.SortDirection,
	) contract.Order,
) (*QueryBuilder[Model, Mutator], error) {
	orderBy = strings.TrimSpace(orderBy)
	if orderBy == "" {
		return qb, nil
	}

	parsedOrders, err := aip132.ParseOrderBy(orderBy)
	if err != nil {
		return nil, fmt.Errorf("invalid order_by syntax -> %w", err)
	}

	for _, orderClause := range parsedOrders {
		fieldPath := orderClause.FieldPath.String()

		orderBuilderFunc, ok := mapColumnNameToOrderFunc[fieldPath]
		if !ok {
			return nil, fmt.Errorf(
				"sorting by field %q is not supported",
				fieldPath,
			)
		}

		direction := contract.SortAsc
		if orderClause.Descending {
			direction = contract.SortDesc
		}

		qb.orders = append(qb.orders, orderBuilderFunc(direction))
	}

	return qb, nil
}

// With appends relation contracts bound to this exact Model type.
func (qb *QueryBuilder[Model, Mutator]) With(
	relations ...contract.Relation[Model],
) *QueryBuilder[Model, Mutator] {
	qb.relations = append(qb.relations, relations...)
	return qb
}

// Distinct marks the query to select unique records.
func (qb *QueryBuilder[Model, Mutator]) Distinct() *QueryBuilder[
	Model,
	Mutator,
] {
	qb.distinct = true
	return qb
}

// GroupBy adds columns to group the query results by.
func (qb *QueryBuilder[Model, Mutator]) GroupBy(
	columns ...string,
) *QueryBuilder[Model, Mutator] {
	qb.groups = append(qb.groups, columns...)
	return qb
}

// Having appends a condition for filtering grouped results.
func (qb *QueryBuilder[Model, Mutator]) Having(
	w contract.Where,
) *QueryBuilder[Model, Mutator] {
	qb.havings = append(qb.havings, w)
	return qb
}

// Get executes the query and returns all matching records.
func (qb *QueryBuilder[Model, Mutator]) Get(
	ctx context.Context,
) ([]Model, error) {
	sql, bindings := qb.compiler.CompileSelect(qb)
	rows, err := qb.db.Query(ctx, sql, bindings...)
	if err != nil {
		return nil, fmt.Errorf(
			"failed executing query -> %w",
			err,
		)
	}
	defer rows.Close()

	cols := qb.GetSelectedColumns()
	var results []Model
	for rows.Next() {
		var m Model
		pointers := qb.hydrate(&m, cols)
		err = rows.Scan(pointers...)
		if err != nil {
			return nil, fmt.Errorf(
				"failed scanning database row -> %w",
				err,
			)
		}
		results = append(results, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf(
			"row iteration error -> %w",
			err,
		)
	}

	if len(qb.relations) > 0 {
		for _, rel := range qb.relations {
			err = rel.EagerLoad(ctx, qb.db, results)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to eager load relation -> %w",
					err,
				)
			}
		}
	}

	return results, nil
}

// First executes the query and returns the first matching record.
func (qb *QueryBuilder[Model, Mutator]) First(
	ctx context.Context,
) (Model, error) {
	cloned := qb.Clone()
	cloned.Limit(1)
	results, err := cloned.Get(ctx)
	if err != nil {
		var zero Model
		return zero, err
	}
	if len(results) == 0 {
		var zero Model
		return zero, errors.New("no matching record found")
	}
	return results[0], nil
}

// Count executes the query and returns the count of matching records.
func (qb *QueryBuilder[Model, Mutator]) Count(
	ctx context.Context,
) (uint64, error) {
	cloned := qb.Clone()
	cloned.orders = nil
	cloned.aggregate = &contract.AggregateState{
		Function: contract.AggCount,
		Column:   "*",
	}
	sql, bindings := qb.compiler.CompileSelect(cloned)
	var count uint64
	row := qb.db.QueryRow(ctx, sql, bindings...)
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf(
			"failed executing count scan -> %w",
			err,
		)
	}
	return count, nil
}

// Exists checks whether any records match the query conditions.
func (qb *QueryBuilder[Model, Mutator]) Exists(
	ctx context.Context,
) (bool, error) {
	count, err := qb.Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Insert inserts a single mutator record into the database.
func (qb *QueryBuilder[Model, Mutator]) Insert(
	ctx context.Context,
	mutator Mutator,
) error {
	return qb.InsertMany(ctx, []Mutator{mutator})
}

// InsertMany inserts multiple mutator records into the database.
func (qb *QueryBuilder[Model, Mutator]) InsertMany(
	ctx context.Context,
	mutators []Mutator,
) error {
	values := make([][]contract.ColumnValue, len(mutators))
	for i := range mutators {
		values[i] = qb.mutatorToColumnValues(&mutators[i])
	}
	sql, bindings := qb.compiler.CompileInsert(qb, values)
	_, err := qb.db.Exec(ctx, sql, bindings...)
	if err != nil {
		return fmt.Errorf(
			"failed executing bulk insert -> %w",
			err,
		)
	}
	return nil
}

// Upsert inserts or updates a single mutator record.
func (qb *QueryBuilder[Model, Mutator]) Upsert(
	ctx context.Context,
	mutator Mutator,
	conflictColumns ...string,
) error {
	return qb.UpsertMany(ctx, []Mutator{mutator}, conflictColumns...)
}

// UpsertMany inserts or updates multiple mutator records.
func (qb *QueryBuilder[Model, Mutator]) UpsertMany(
	ctx context.Context,
	mutators []Mutator,
	conflictColumns ...string,
) error {
	values := make([][]contract.ColumnValue, len(mutators))
	for i := range mutators {
		values[i] = qb.mutatorToColumnValues(&mutators[i])
	}
	sql, bindings := qb.compiler.CompileUpsert(
		qb,
		values,
		conflictColumns,
	)
	_, err := qb.db.Exec(ctx, sql, bindings...)
	if err != nil {
		return fmt.Errorf(
			"failed executing bulk upsert -> %w",
			err,
		)
	}
	return nil
}

// Update updates records matching the query conditions.
func (qb *QueryBuilder[Model, Mutator]) Update(
	ctx context.Context,
	mutator Mutator,
) error {
	values := qb.mutatorToColumnValues(&mutator)
	sql, bindings := qb.compiler.CompileUpdate(qb, values)
	_, err := qb.db.Exec(ctx, sql, bindings...)
	if err != nil {
		return fmt.Errorf(
			"failed executing update -> %w",
			err,
		)
	}
	return nil
}

// Delete deletes records matching the query conditions.
func (qb *QueryBuilder[Model, Mutator]) Delete(
	ctx context.Context,
) error {
	sql, bindings := qb.compiler.CompileDelete(qb)
	_, err := qb.db.Exec(ctx, sql, bindings...)
	if err != nil {
		return fmt.Errorf(
			"failed executing delete -> %w",
			err,
		)
	}
	return nil
}

// InsertReturning inserts a single record and returns the hydrated model.
func (qb *QueryBuilder[Model, Mutator]) InsertReturning(
	ctx context.Context,
	mutator Mutator,
) (Model, error) {
	results, err := qb.InsertReturningMany(ctx, []Mutator{mutator})
	if err != nil {
		var zero Model
		return zero, err
	}
	if len(results) == 0 {
		var zero Model
		return zero, errors.New(
			"insert returning failed -> no record returned",
		)
	}
	return results[0], nil
}

// InsertReturningMany inserts multiple records and returns hydrated models.
func (qb *QueryBuilder[Model, Mutator]) InsertReturningMany(
	ctxQuery context.Context,
	mutators []Mutator,
) ([]Model, error) {
	values := make([][]contract.ColumnValue, len(mutators))
	for i := range mutators {
		values[i] = qb.mutatorToColumnValues(&mutators[i])
	}
	cols := qb.GetDefaultColumns()
	sql, bindings := qb.compiler.CompileInsertReturning(
		qb,
		values,
		cols,
	)
	rows, err := qb.db.Query(ctxQuery, sql, bindings...)
	if err != nil {
		return nil, fmt.Errorf(
			"failed executing insert returning -> %w",
			err,
		)
	}
	defer rows.Close()

	results := make([]Model, 0, len(mutators))
	for rows.Next() {
		var m Model
		pointers := qb.hydrate(&m, cols)
		err = rows.Scan(pointers...)
		if err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("row iteration error -> %w", err)
	}
	return results, nil
}

// UpsertReturning upserts a single record and returns the hydrated model.
func (qb *QueryBuilder[Model, Mutator]) UpsertReturning(
	ctx context.Context,
	mutator Mutator,
	conflictColumns ...string,
) (Model, error) {
	results, err := qb.UpsertReturningMany(
		ctx,
		[]Mutator{mutator},
		conflictColumns...,
	)
	if err != nil {
		var zero Model
		return zero, err
	}
	if len(results) == 0 {
		var zero Model
		return zero, errors.New(
			"upsert returning failed -> no record returned",
		)
	}
	return results[0], nil
}

// UpsertReturningMany upserts multiple records and returns hydrated models.
func (qb *QueryBuilder[Model, Mutator]) UpsertReturningMany(
	ctxQuery context.Context,
	mutators []Mutator,
	conflictColumns ...string,
) ([]Model, error) {
	values := make([][]contract.ColumnValue, len(mutators))
	for i := range mutators {
		values[i] = qb.mutatorToColumnValues(&mutators[i])
	}
	cols := qb.GetDefaultColumns()
	sql, bindings := qb.compiler.CompileUpsert(
		qb,
		values,
		conflictColumns,
	)
	sql += " RETURNING " + strings.Join(cols, ", ")
	rows, err := qb.db.Query(ctxQuery, sql, bindings...)
	if err != nil {
		return nil, fmt.Errorf(
			"failed executing upsert returning -> %w",
			err,
		)
	}
	defer rows.Close()

	results := make([]Model, 0, len(mutators))
	for rows.Next() {
		var m Model
		pointers := qb.hydrate(&m, cols)
		err = rows.Scan(pointers...)
		if err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("row iteration error -> %w", err)
	}
	return results, nil
}

// UpdateReturning updates records and returns the hydrated models.
func (qb *QueryBuilder[Model, Mutator]) UpdateReturning(
	ctxQuery context.Context,
	mutator Mutator,
) ([]Model, error) {
	values := qb.mutatorToColumnValues(&mutator)
	cols := qb.GetDefaultColumns()
	sql, bindings := qb.compiler.CompileUpdateReturning(
		qb,
		values,
		cols,
	)
	rows, err := qb.db.Query(ctxQuery, sql, bindings...)
	if err != nil {
		return nil, fmt.Errorf(
			"failed executing update returning -> %w",
			err,
		)
	}
	defer rows.Close()

	var results []Model
	for rows.Next() {
		var m Model
		pointers := qb.hydrate(&m, cols)
		err = rows.Scan(pointers...)
		if err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("row iteration error -> %w", err)
	}
	return results, nil
}

// DeleteReturning deletes records and returns the hydrated models.
func (qb *QueryBuilder[Model, Mutator]) DeleteReturning(
	ctxQuery context.Context,
) ([]Model, error) {
	cols := qb.GetDefaultColumns()
	sql, bindings := qb.compiler.CompileDeleteReturning(
		qb,
		cols,
	)
	rows, err := qb.db.Query(ctxQuery, sql, bindings...)
	if err != nil {
		return nil, fmt.Errorf(
			"failed executing delete returning -> %w",
			err,
		)
	}
	defer rows.Close()

	var results []Model
	for rows.Next() {
		var m Model
		pointers := qb.hydrate(&m, cols)
		err = rows.Scan(pointers...)
		if err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("row iteration error -> %w", err)
	}
	return results, nil
}

// Paginate retrieves a slice of models within the limit and offset.
func (qb *QueryBuilder[Model, Mutator]) Paginate(
	ctx context.Context,
	limit uint64,
	offset uint64,
) (contract.PaginateResult[Model], error) {
	cloned := qb.Clone()
	cloned.Limit(limit + 1).Offset(offset)
	items, err := cloned.Get(ctx)
	if err != nil {
		return contract.PaginateResult[Model]{}, err
	}

	hasMore := false
	if len(items) > int(limit) {
		hasMore = true
		items = items[:limit]
	}
	return contract.PaginateResult[Model]{
		Items:   items,
		HasMore: hasMore,
	}, nil
}

// PaginateAip158 implements AIP-158 pagination for the query.
func (qb *QueryBuilder[Model, Mutator]) PaginateAip158(
	ctx context.Context,
	req pagination.Request,
	defaultPageSize int32,
) (contract.PaginateAip158Result[Model], error) {
	pageToken, err := pagination.ParsePageToken(req)
	if err != nil {
		return contract.PaginateAip158Result[Model]{}, fmt.Errorf(
			"failed to parse page token -> %w",
			err,
		)
	}

	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	cloned := qb.Clone()
	cloned.Limit(uint64(pageSize) + 1).Offset(uint64(pageToken.Offset))
	items, err := cloned.Get(ctx)
	if err != nil {
		return contract.PaginateAip158Result[Model]{}, err
	}

	var nextPageToken string
	if len(items) > int(pageSize) {
		items = items[:pageSize]
		nextPageToken = pageToken.Next(req).String()
	}

	return contract.PaginateAip158Result[Model]{
		Items:         items,
		NextPageToken: nextPageToken,
	}, nil
}

// PaginateWithCount retrieves a slice of models and total count.
func (qb *QueryBuilder[Model, Mutator]) PaginateWithCount(
	ctx context.Context,
	limit uint64,
	offset uint64,
) (contract.PaginateCountResult[Model], error) {
	total, err := qb.Count(ctx)
	if err != nil {
		return contract.PaginateCountResult[Model]{}, err
	}

	cloned := qb.Clone()
	cloned.Limit(limit).Offset(offset)
	items, err := cloned.Get(ctx)
	if err != nil {
		return contract.PaginateCountResult[Model]{}, err
	}

	return contract.PaginateCountResult[Model]{
		Items:      items,
		TotalCount: int64(total),
	}, nil
}

// Min calculates the minimum value of a column.
func (qb *QueryBuilder[Model, Mutator]) Min[T any](
	ctx context.Context,
	column interface {
		contract.Column[Model]
		BindType() T
	},
) (T, error) {
	cloned := qb.Clone()
	cloned.aggregate = &contract.AggregateState{
		Function: contract.AggMin,
		Column:   column.ColumnName(),
	}
	sql, bindings := qb.compiler.CompileSelect(cloned)
	var res T
	row := qb.db.QueryRow(ctx, sql, bindings...)
	err := row.Scan(&res)
	if err != nil {
		return res, fmt.Errorf(
			"failed executing min aggregate scan -> %w",
			err,
		)
	}
	return res, nil
}

// Max calculates the maximum value of a column.
func (qb *QueryBuilder[Model, Mutator]) Max[T any](
	ctx context.Context,
	column interface {
		contract.Column[Model]
		BindType() T
	},
) (T, error) {
	cloned := qb.Clone()
	cloned.aggregate = &contract.AggregateState{
		Function: contract.AggMax,
		Column:   column.ColumnName(),
	}
	sql, bindings := qb.compiler.CompileSelect(cloned)
	var res T
	row := qb.db.QueryRow(ctx, sql, bindings...)
	err := row.Scan(&res)
	if err != nil {
		return res, fmt.Errorf(
			"failed executing max aggregate scan -> %w",
			err,
		)
	}
	return res, nil
}

// Avg calculates the average value of a column.
func (qb *QueryBuilder[Model, Mutator]) Avg[T any](
	ctx context.Context,
	column interface {
		contract.Column[Model]
		BindType() T
	},
) (T, error) {
	cloned := qb.Clone()
	cloned.aggregate = &contract.AggregateState{
		Function: contract.AggAvg,
		Column:   column.ColumnName(),
	}
	sql, bindings := qb.compiler.CompileSelect(cloned)
	var res T
	row := qb.db.QueryRow(ctx, sql, bindings...)
	err := row.Scan(&res)
	if err != nil {
		return res, fmt.Errorf(
			"failed executing avg aggregate scan -> %w",
			err,
		)
	}
	return res, nil
}

// Sum calculates the sum value of a column.
func (qb *QueryBuilder[Model, Mutator]) Sum[T any](
	ctx context.Context,
	column interface {
		contract.Column[Model]
		BindType() T
	},
) (T, error) {
	cloned := qb.Clone()
	cloned.aggregate = &contract.AggregateState{
		Function: contract.AggSum,
		Column:   column.ColumnName(),
	}
	sql, bindings := qb.compiler.CompileSelect(cloned)
	var res T
	row := qb.db.QueryRow(ctx, sql, bindings...)
	err := row.Scan(&res)
	if err != nil {
		return res, fmt.Errorf(
			"failed executing sum aggregate scan -> %w",
			err,
		)
	}
	return res, nil
}
