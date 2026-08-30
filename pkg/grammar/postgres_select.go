package grammar

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// CompileSelect compiles a select QueryStateProvider into SELECT SQL.
func (g *PostgresGrammar) CompileSelect(
	state contract.QueryStateProvider,
) (string, []any) {
	tracker := bindingsTracker{values: make([]any, 0)}

	sql, tracker := g.compileSelectWithTracker(state, tracker)
	return sql, tracker.values
}

// compileSelectWithTracker compiles a select query with the bindings
// tracker.
func (g *PostgresGrammar) compileSelectWithTracker(
	state contract.QueryStateProvider,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	components := []func(
		g *PostgresGrammar,
		state contract.QueryStateProvider,
		tracker bindingsTracker,
	) (string, bindingsTracker){
		(*PostgresGrammar).compileAggregate,
		(*PostgresGrammar).compileColumns,
		(*PostgresGrammar).compileFrom,
		(*PostgresGrammar).compileWheres,
		(*PostgresGrammar).compileGroups,
		(*PostgresGrammar).compileHavings,
		(*PostgresGrammar).compileOrders,
		(*PostgresGrammar).compileLimit,
		(*PostgresGrammar).compileOffset,
	}

	var sqlParts []string
	for _, componentFn := range components {
		var part string
		part, tracker = componentFn(g, state, tracker)
		if part == "" {
			continue
		}
		sqlParts = append(sqlParts, part)
	}

	return strings.Join(sqlParts, " "), tracker
}

// compileAggregate compiles select aggregate queries.
func (g *PostgresGrammar) compileAggregate(
	state contract.QueryStateProvider,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	agg := state.GetAggregate()
	if agg == nil {
		return "", tracker
	}

	if state.IsDistinct() {
		return fmt.Sprintf(
			"SELECT %s(DISTINCT %s) AS aggregate",
			agg.Function,
			sanitizeColumn(agg.Column),
		), tracker
	}
	return fmt.Sprintf(
		"SELECT %s(%s) AS aggregate",
		agg.Function,
		sanitizeColumn(agg.Column),
	), tracker
}

// compileColumns compiles select column projection lists.
func (g *PostgresGrammar) compileColumns(
	state contract.QueryStateProvider,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	// If aggregate compilation already handled selection, skip columns
	if state.GetAggregate() != nil {
		return "", tracker
	}

	cols := state.GetSelectedColumns()
	var columnsSql string
	switch len(cols) {
	case 0:
		columnsSql = "*"
	default:
		wrapped := make([]string, len(cols))
		for i, col := range cols {
			wrapped[i] = sanitizeColumn(col)
		}
		columnsSql = strings.Join(wrapped, ", ")
	}

	if state.IsDistinct() {
		return fmt.Sprintf("SELECT DISTINCT %s", columnsSql), tracker
	}
	return fmt.Sprintf("SELECT %s", columnsSql), tracker
}

// compileFrom compiles the FROM table segment.
func (g *PostgresGrammar) compileFrom(
	state contract.QueryStateProvider,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	sanitized := pgx.Identifier{
		state.GetSchemaName(),
		state.GetTableName(),
	}.Sanitize()
	return fmt.Sprintf("FROM %s", sanitized), tracker
}

// compileGroups compiles query GROUP BY columns.
func (g *PostgresGrammar) compileGroups(
	state contract.QueryStateProvider,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	groups := state.GetGroups()
	if len(groups) == 0 {
		return "", tracker
	}
	wrapped := make([]string, len(groups))
	for i, col := range groups {
		wrapped[i] = sanitizeColumn(col)
	}
	return fmt.Sprintf("GROUP BY %s", strings.Join(wrapped, ", ")), tracker
}

// compileLimit compiles query LIMIT boundaries.
func (g *PostgresGrammar) compileLimit(
	state contract.QueryStateProvider,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	lim, ok := state.GetLimit()
	if !ok {
		return "", tracker
	}
	return fmt.Sprintf("LIMIT %d", lim), tracker
}

// compileOffset compiles query OFFSET skips.
func (g *PostgresGrammar) compileOffset(
	state contract.QueryStateProvider,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	off, ok := state.GetOffset()
	if !ok {
		return "", tracker
	}
	return fmt.Sprintf("OFFSET %d", off), tracker
}
