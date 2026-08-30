package grammar

import (
	"fmt"
	"strings"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// compileWheres compiles a query state's WHERE clauses.
func (g *PostgresGrammar) compileWheres(
	state contract.QueryStateProvider,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	var wheresSql string
	wheresSql, tracker = g.compileWheresRaw(state.GetWheres(), tracker)
	if wheresSql == "" {
		return "", tracker
	}
	return fmt.Sprintf("WHERE %s", wheresSql), tracker
}

// compileHavings compiles a query state's HAVING clauses.
func (g *PostgresGrammar) compileHavings(
	state contract.QueryStateProvider,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	var havingsSql string
	havingsSql, tracker = g.compileWheresRaw(state.GetHavings(), tracker)
	if havingsSql == "" {
		return "", tracker
	}
	return fmt.Sprintf("HAVING %s", havingsSql), tracker
}

// compileWheresRaw transforms a slice of generic Where interfaces into SQL.
func (g *PostgresGrammar) compileWheresRaw(
	wheres []contract.Where,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	if len(wheres) == 0 {
		return "", tracker
	}

	var sb strings.Builder

	for _, w := range wheres {
		if w == nil {
			continue
		}

		var conditionSql string
		var boolean string

		switch cond := w.(type) {
		case contract.JsonWhere:
			conditionSql, tracker = g.compileJsonWhere(
				cond,
				tracker,
			)
			boolean = string(cond.GetBoolean())
		case contract.SubQueryWhere:
			conditionSql, tracker = g.compileSubQueryWhere(
				cond,
				tracker,
			)
			boolean = string(cond.GetBoolean())
		case contract.DateWhere:
			conditionSql, tracker = g.compileDateWhere(
				cond,
				tracker,
			)
			boolean = string(cond.GetBoolean())
		case contract.ColumnWhere:
			conditionSql, tracker = g.compileColumnWhere(
				cond,
				tracker,
			)
			boolean = string(cond.GetBoolean())
		case contract.BetweenWhere:
			conditionSql, tracker = g.compileBetweenWhere(
				cond,
				tracker,
			)
			boolean = string(cond.GetBoolean())
		case contract.InWhere:
			conditionSql, tracker = g.compileInWhere(
				cond,
				tracker,
			)
			boolean = string(cond.GetBoolean())
		case contract.BasicWhere:
			conditionSql, tracker = g.compileBasicWhere(
				cond,
				tracker,
			)
			boolean = string(cond.GetBoolean())
		case contract.ExistsWhere:
			conditionSql, tracker = g.compileExistsWhere(
				cond,
				tracker,
			)
			boolean = string(cond.GetBoolean())
		case contract.NullWhere:
			conditionSql, tracker = g.compileNullWhere(
				cond,
				tracker,
			)
			boolean = string(cond.GetBoolean())
		case contract.NestedWhere:
			conditionSql, tracker = g.compileNestedWhere(
				cond,
				tracker,
			)
			boolean = string(cond.GetBoolean())
		case contract.RawWhere:
			conditionSql, tracker = g.compileRawWhere(
				cond,
				tracker,
			)
			boolean = string(cond.GetBoolean())
		default:
			panic(fmt.Sprintf("unknown Where condition type: %T", w))
		}

		if sb.Len() > 0 {
			sb.WriteString(fmt.Sprintf(" %s ", boolean))
		}
		sb.WriteString(conditionSql)
	}

	return sb.String(), tracker
}

// compileBasicWhere compiles standard column-value conditions.
func (g *PostgresGrammar) compileBasicWhere(
	w contract.BasicWhere,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	var placeholder string
	placeholder, tracker = tracker.Bind(w.GetValue())
	sql := fmt.Sprintf(
		"%s %s %s",
		sanitizeColumn(w.GetColumn()),
		string(w.GetOperator()),
		placeholder,
	)
	return sql, tracker
}

// compileNullWhere compiles null check conditions.
func (g *PostgresGrammar) compileNullWhere(
	w contract.NullWhere,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	op := "IS NULL"
	if w.IsNot() {
		op = "IS NOT NULL"
	}
	sql := fmt.Sprintf("%s %s", sanitizeColumn(w.GetColumn()), op)
	return sql, tracker
}

// compileInWhere compiles membership conditions.
func (g *PostgresGrammar) compileInWhere(
	w contract.InWhere,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	vals := w.GetValues()
	if len(vals) == 0 {
		if w.IsNot() {
			return "1 = 1", tracker
		}
		return "1 = 0", tracker
	}

	placeholders := make([]string, len(vals))
	for i, val := range vals {
		placeholders[i], tracker = tracker.Bind(val)
	}

	op := string(contract.OpIn)
	if w.IsNot() {
		op = string(contract.OpNotIn)
	}

	sql := fmt.Sprintf(
		"%s %s (%s)",
		sanitizeColumn(w.GetColumn()),
		op,
		strings.Join(placeholders, ", "),
	)
	return sql, tracker
}

// compileBetweenWhere compiles range-boundary conditions.
func (g *PostgresGrammar) compileBetweenWhere(
	w contract.BetweenWhere,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	op := string(contract.OpBetween)
	if w.IsNot() {
		op = string(contract.OpNotBetween)
	}
	var pMin string
	var pMax string
	pMin, tracker = tracker.Bind(w.GetMin())
	pMax, tracker = tracker.Bind(w.GetMax())
	sql := fmt.Sprintf(
		"%s %s %s AND %s",
		sanitizeColumn(w.GetColumn()),
		op,
		pMin,
		pMax,
	)
	return sql, tracker
}

// compileNestedWhere compiles logical sub-groups enclosed in parentheses.
func (g *PostgresGrammar) compileNestedWhere(
	w contract.NestedWhere,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	var sql string
	sql, tracker = g.compileWheresRaw(
		w.GetQuery().GetWheres(),
		tracker,
	)
	prefix := ""
	if w.IsNegated() {
		prefix = "NOT "
	}
	return fmt.Sprintf("%s(%s)", prefix, sql), tracker
}

// compileRawWhere compiles unsafely specified raw SQL strings.
func (g *PostgresGrammar) compileRawWhere(
	w contract.RawWhere,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	sql := w.GetSql()
	bindings := w.GetBindings()

	for _, bindVal := range bindings {
		var placeholder string
		placeholder, tracker = tracker.Bind(bindVal)
		sql = strings.Replace(sql, "?", placeholder, 1)
	}

	return sql, tracker
}

// compileColumnWhere compiles comparisons of columns with other columns.
func (g *PostgresGrammar) compileColumnWhere(
	w contract.ColumnWhere,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	sql := fmt.Sprintf(
		"%s %s %s",
		sanitizeColumn(w.GetFirst()),
		string(w.GetOperator()),
		sanitizeColumn(w.GetSecond()),
	)
	return sql, tracker
}

// compileDateWhere compiles date-specific operations.
func (g *PostgresGrammar) compileDateWhere(
	w contract.DateWhere,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	column := sanitizeColumn(w.GetColumn())
	var placeholder string
	placeholder, tracker = tracker.Bind(w.GetValue())
	operator := string(w.GetOperator())

	switch strings.ToLower(w.GetType()) {
	case "date":
		return fmt.Sprintf(
			"%s::date %s %s",
			column,
			operator,
			placeholder,
		), tracker
	case "year":
		return fmt.Sprintf(
			"EXTRACT(YEAR FROM %s) %s %s",
			column,
			operator,
			placeholder,
		), tracker
	case "month":
		return fmt.Sprintf(
			"EXTRACT(MONTH FROM %s) %s %s",
			column,
			operator,
			placeholder,
		), tracker
	case "day":
		return fmt.Sprintf(
			"EXTRACT(DAY FROM %s) %s %s",
			column,
			operator,
			placeholder,
		), tracker
	case "time":
		return fmt.Sprintf(
			"%s::time %s %s",
			column,
			operator,
			placeholder,
		), tracker
	default:
		panic(fmt.Sprintf(
			"unsupported DateWhere extract type: %s",
			w.GetType(),
		))
	}
}

// compileExistsWhere compiles existence assertions.
func (g *PostgresGrammar) compileExistsWhere(
	w contract.ExistsWhere,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	var subQuerySql string
	subQuerySql, tracker = g.compileSelectWithTracker(
		w.GetQuery(),
		tracker,
	)
	op := "EXISTS"
	if w.IsNot() {
		op = "NOT EXISTS"
	}
	return fmt.Sprintf("%s (%s)", op, subQuerySql), tracker
}

// compileSubQueryWhere compiles sub-query scalar checks.
func (g *PostgresGrammar) compileSubQueryWhere(
	w contract.SubQueryWhere,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	var subQuerySql string
	subQuerySql, tracker = g.compileSelectWithTracker(
		w.GetQuery(),
		tracker,
	)
	sql := fmt.Sprintf(
		"%s %s (%s)",
		sanitizeColumn(w.GetColumn()),
		string(w.GetOperator()),
		subQuerySql,
	)
	return sql, tracker
}

// compileJsonWhere compiles JSON path-key extract value conditions.
func (g *PostgresGrammar) compileJsonWhere(
	w contract.JsonWhere,
	tracker bindingsTracker,
) (string, bindingsTracker) {
	var pKey string
	var pVal string
	pKey, tracker = tracker.Bind(w.GetKey())
	pVal, tracker = tracker.Bind(w.GetValue())
	sql := fmt.Sprintf(
		"%s %s %s %s %s",
		sanitizeColumn(w.GetColumn()),
		string(contract.OpJsonGetTextField),
		pKey,
		string(w.GetOperator()),
		pVal,
	)
	return sql, tracker
}
