package orm

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// mutatorToColumnValues converts a mutator to column values.
func (qb *QueryBuilder[Model, Mutator]) mutatorToColumnValues(
	m *Mutator,
) []contract.ColumnValue {
	cols, vals := qb.dehydrate(m)
	res := make([]contract.ColumnValue, len(cols))
	for i := range cols {
		res[i] = contract.ColumnValue{
			Column: cols[i],
			Value:  vals[i],
		}
	}
	return res
}
