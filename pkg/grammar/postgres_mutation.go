package grammar

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// CompileInsert compiles an INSERT statement for multiple rows of column
// values using high-performance Postgres UNNEST array functions.
func (g *PostgresGrammar) CompileInsert(
	state contract.QueryStateProvider,
	values [][]contract.ColumnValue,
) (string, []any) {
	if len(values) == 0 {
		return "", nil
	}
	tracker := bindingsTracker{values: make([]any, 0)}

	schemaName := state.GetSchemaName()
	tableName := state.GetTableName()
	table := pgx.Identifier{schemaName, tableName}.Sanitize()

	var rawColumns []string
	mapColumnToIndex := make(map[string]int)
	for _, row := range values {
		for _, colVal := range row {
			_, ok := mapColumnToIndex[colVal.Column]
			if !ok {
				mapColumnToIndex[colVal.Column] = len(rawColumns)
				rawColumns = append(rawColumns, colVal.Column)
			}
		}
	}

	columns := make([]string, len(rawColumns))
	for i, col := range rawColumns {
		columns[i] = sanitizeColumn(col)
	}

	// Vertically slice values into per-column arrays using a single
	// contiguous backing slice.
	backing := make([]any, len(columns)*len(values))
	columnArrays := make([][]any, len(columns))
	for i := range columns {
		columnArrays[i] = backing[i*len(values) : (i+1)*len(values)]
	}

	for rowIndex, row := range values {
		for _, colVal := range row {
			colIndex := mapColumnToIndex[colVal.Column]
			columnArrays[colIndex][rowIndex] = colVal.Value
		}
	}

	// Bind each column's array as a single placeholder with explicit
	// Postgres casting.
	placeholders := make([]string, len(columns))
	selectCols := make([]string, len(columns))
	aliasCols := make([]string, len(columns))
	for i, colName := range rawColumns {
		suffix, typedArray, isArray := state.GetColumnCastAndTypedSlice(
			colName,
			columnArrays[i],
		)
		ph, tk := tracker.Bind(typedArray)
		tracker = tk
		placeholders[i] = ph + suffix

		colIdx := i + 1
		aliasCols[i] = fmt.Sprintf("unnested_column_%d", colIdx)
		switch isArray {
		case true:
			selectCols[i] = fmt.Sprintf(
				"unnested.unnested_column_%d%s",
				colIdx,
				suffix,
			)
		case false:
			selectCols[i] = fmt.Sprintf(
				"unnested.unnested_column_%d",
				colIdx,
			)
		}
	}

	sql := fmt.Sprintf(
		"INSERT INTO %s (%s) SELECT %s FROM UNNEST(%s) AS unnested(%s)",
		table,
		strings.Join(columns, ", "),
		strings.Join(selectCols, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(aliasCols, ", "),
	)
	return sql, tracker.values
}

// CompileInsertReturning compiles an INSERT statement and appends a
// RETURNING clause.
func (g *PostgresGrammar) CompileInsertReturning(
	state contract.QueryStateProvider,
	values [][]contract.ColumnValue,
	returning []string,
) (string, []any) {
	sql, bindings := g.CompileInsert(state, values)
	if sql == "" {
		return "", nil
	}
	if len(returning) == 0 {
		return sql, bindings
	}
	wrapped := make([]string, len(returning))
	for i, col := range returning {
		wrapped[i] = sanitizeColumn(col)
	}
	sql = fmt.Sprintf("%s RETURNING %s", sql, strings.Join(wrapped, ", "))
	return sql, bindings
}

// CompileUpdate compiles an UPDATE statement for a slice of assignment
// values.
func (g *PostgresGrammar) CompileUpdate(
	state contract.QueryStateProvider,
	values []contract.ColumnValue,
) (string, []any) {
	if len(values) == 0 {
		return "", nil
	}
	tracker := bindingsTracker{values: make([]any, 0)}

	schemaName := state.GetSchemaName()
	tableName := state.GetTableName()
	table := pgx.Identifier{schemaName, tableName}.Sanitize()

	assignments := make([]string, len(values))
	for i, val := range values {
		var placeholder string
		placeholder, tracker = tracker.Bind(val.Value)
		assignments[i] = fmt.Sprintf(
			"%s = %s",
			sanitizeColumn(val.Column),
			placeholder,
		)
	}

	sql := fmt.Sprintf("UPDATE %s SET %s", table, strings.Join(assignments, ", "))

	var wheresSql string
	wheresSql, tracker = g.compileWheresRaw(
		state.GetWheres(),
		tracker,
	)
	if wheresSql != "" {
		sql = fmt.Sprintf("%s WHERE %s", sql, wheresSql)
	}

	return sql, tracker.values
}

// CompileUpdateReturning compiles an UPDATE statement and appends a
// RETURNING clause.
func (g *PostgresGrammar) CompileUpdateReturning(
	state contract.QueryStateProvider,
	values []contract.ColumnValue,
	returning []string,
) (string, []any) {
	sql, bindings := g.CompileUpdate(state, values)
	if sql == "" {
		return "", nil
	}
	if len(returning) == 0 {
		return sql, bindings
	}
	wrapped := make([]string, len(returning))
	for i, col := range returning {
		wrapped[i] = sanitizeColumn(col)
	}
	sql = fmt.Sprintf("%s RETURNING %s", sql, strings.Join(wrapped, ", "))
	return sql, bindings
}

// CompileDelete compiles a DELETE statement.
func (g *PostgresGrammar) CompileDelete(
	state contract.QueryStateProvider,
) (string, []any) {
	tracker := bindingsTracker{values: make([]any, 0)}

	schemaName := state.GetSchemaName()
	tableName := state.GetTableName()
	table := pgx.Identifier{schemaName, tableName}.Sanitize()

	sql := fmt.Sprintf("DELETE FROM %s", table)

	var wheresSql string
	wheresSql, tracker = g.compileWheresRaw(state.GetWheres(), tracker)
	if wheresSql != "" {
		sql = fmt.Sprintf("%s WHERE %s", sql, wheresSql)
	}

	return sql, tracker.values
}

// CompileDeleteReturning compiles a DELETE statement and appends a
// RETURNING clause.
func (g *PostgresGrammar) CompileDeleteReturning(
	state contract.QueryStateProvider,
	returning []string,
) (string, []any) {
	sql, bindings := g.CompileDelete(state)
	if len(returning) == 0 {
		return sql, bindings
	}
	wrapped := make([]string, len(returning))
	for i, col := range returning {
		wrapped[i] = sanitizeColumn(col)
	}
	sql = fmt.Sprintf("%s RETURNING %s", sql, strings.Join(wrapped, ", "))
	return sql, bindings
}

// CompileUpsert compiles an INSERT statement with conflict-triggered ON
// CONFLICT update logic.
func (g *PostgresGrammar) CompileUpsert(
	state contract.QueryStateProvider,
	values [][]contract.ColumnValue,
	conflictColumns []string,
) (string, []any) {
	sql, bindings := g.CompileInsert(state, values)
	if sql == "" {
		return "", nil
	}

	wrappedConflicts := make([]string, len(conflictColumns))
	mapConflictColumnToIsConflict := make(map[string]bool)
	for i, col := range conflictColumns {
		wrappedConflicts[i] = sanitizeColumn(col)
		mapConflictColumnToIsConflict[col] = true
	}

	var rawColumns []string
	mapColumnToIsIncluded := make(map[string]bool)
	for _, row := range values {
		for _, colVal := range row {
			_, hasCol := mapColumnToIsIncluded[colVal.Column]
			if !hasCol {
				mapColumnToIsIncluded[colVal.Column] = true
				rawColumns = append(rawColumns, colVal.Column)
			}
		}
	}

	var updates []string
	for _, col := range rawColumns {
		_, isConflict := mapConflictColumnToIsConflict[col]
		if !isConflict {
			wrappedCol := sanitizeColumn(col)
			updates = append(
				updates,
				fmt.Sprintf("%s = EXCLUDED.%s", wrappedCol, wrappedCol),
			)
		}
	}

	upsertSql := fmt.Sprintf(
		"%s ON CONFLICT (%s) DO ",
		sql,
		strings.Join(wrappedConflicts, ", "),
	)

	if len(updates) > 0 {
		updatesStr := strings.Join(updates, ", ")
		upsertSql = fmt.Sprintf("%sUPDATE SET %s", upsertSql, updatesStr)
	} else {
		upsertSql = fmt.Sprintf("%sNOTHING", upsertSql)
	}

	return upsertSql, bindings
}
