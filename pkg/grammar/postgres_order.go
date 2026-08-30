package grammar

import (
	"fmt"
	"strings"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// compileOrders compiles query ORDER BY clauses into SQL.
func (g *PostgresGrammar) compileOrders(
	state contract.QueryStateProvider,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	orders := state.GetOrders()
	if len(orders) == 0 {
		return "", tracker
	}

	compiledParts := make([]string, len(orders))
	for i, orderClause := range orders {
		compiledParts[i], tracker = g.compileOrder(orderClause, tracker)
	}

	orderStr := strings.Join(compiledParts, ", ")
	return fmt.Sprintf("ORDER BY %s", orderStr), tracker
}

// ==========================================================================
// Specialized Order Compilers (Immutable)
// ==========================================================================

// compileOrder compiles a single Order clause into SQL.
func (g *PostgresGrammar) compileOrder(
	order contract.Order,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	switch ord := order.(type) {
	case contract.ColumnOrder:
		return g.compileColumnOrder(ord), tracker
	case contract.SubQueryOrder:
		return g.compileSubQueryOrder(ord, tracker)
	case contract.RawOrder:
		return g.compileRawOrder(ord), tracker
	default:
		panic(fmt.Sprintf("unknown Order clause type: %T", order))
	}
}

// compileRawOrder compiles a raw SQL sorting expression.
func (g *PostgresGrammar) compileRawOrder(order contract.RawOrder) string {
	sql := order.GetSql()
	direction := g.compileSortDirection(order.GetDirection(), "")
	if direction != "" {
		sql = fmt.Sprintf("%s %s", sql, direction)
	}
	return sql
}

// compileColumnOrder compiles a standard column sort order.
func (g *PostgresGrammar) compileColumnOrder(
	order contract.ColumnOrder,
) string {
	column := sanitizeColumn(order.GetColumn())
	direction := g.compileSortDirection(
		order.GetDirection(),
		order.GetUsingOperator(),
	)
	nulls := string(order.GetNulls())

	sql := column
	if direction != "" {
		sql = fmt.Sprintf("%s %s", sql, direction)
	}
	if nulls != "" {
		sql = fmt.Sprintf("%s %s", sql, nulls)
	}
	return sql
}

// compileSubQueryOrder compiles a sort order based on a scalar subquery.
func (g *PostgresGrammar) compileSubQueryOrder(
	order contract.SubQueryOrder,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	var subQuerySql string
	subQuerySql, tracker = g.compileSelectWithTracker(
		order.GetQuery(),
		tracker,
	)
	direction := g.compileSortDirection(order.GetDirection(), "")
	nulls := string(order.GetNulls())

	sql := fmt.Sprintf("(%s)", subQuerySql)
	if direction != "" {
		sql = fmt.Sprintf("%s %s", sql, direction)
	}
	if nulls != "" {
		sql = fmt.Sprintf("%s %s", sql, nulls)
	}
	return sql, tracker
}

// compileSortDirection compiles the sorting direction or custom USING
// operator.
func (g *PostgresGrammar) compileSortDirection(
	direction contract.SortDirection,
	customOperator contract.ComparisonOperator,
) string {
	if direction == contract.SortUsing {
		return fmt.Sprintf(string(direction), customOperator)
	}
	return string(direction)
}
