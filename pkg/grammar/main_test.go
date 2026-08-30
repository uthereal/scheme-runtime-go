package grammar

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// mockQueryState is a mock implementation of contract.QueryStateProvider
// for testing.
type mockQueryState struct {
	schemaName     string
	tableName      string
	defaultColumns []string
	columns        []string
	distinct       bool
	aggregate      *contract.AggregateState
	wheres         []contract.Where
	orders         []contract.Order
	groups         []string
	havings        []contract.Where
	limit          *uint64
	offset         *uint64
}

// GetSchemaName returns the schema name of the query state.
func (m mockQueryState) GetSchemaName() string {
	return m.schemaName
}

// GetTableName returns the table name of the query state.
func (m mockQueryState) GetTableName() string {
	return m.tableName
}

// GetDefaultColumns returns the default columns of the query state.
func (m mockQueryState) GetDefaultColumns() []string {
	return m.defaultColumns
}

// GetSelectedColumns returns the selected columns of the query state.
func (m mockQueryState) GetSelectedColumns() []string {
	return m.columns
}

// IsDistinct returns whether the query is distinct.
func (m mockQueryState) IsDistinct() bool {
	return m.distinct
}

// GetAggregate returns the aggregate state of the query state.
func (m mockQueryState) GetAggregate() *contract.AggregateState {
	return m.aggregate
}

// GetWheres returns the where clauses of the query state.
func (m mockQueryState) GetWheres() []contract.Where {
	return m.wheres
}

// GetOrders returns the order clauses of the query state.
func (m mockQueryState) GetOrders() []contract.Order {
	return m.orders
}

// GetGroups returns the group clauses of the query state.
func (m mockQueryState) GetGroups() []string {
	return m.groups
}

// GetHavings returns the having clauses of the query state.
func (m mockQueryState) GetHavings() []contract.Where {
	return m.havings
}

// GetLimit returns the limit value of the query state.
func (m mockQueryState) GetLimit() (uint64, bool) {
	if m.limit == nil {
		return 0, false
	}
	return *m.limit, true
}

// GetOffset returns the offset value of the query state.
func (m mockQueryState) GetOffset() (uint64, bool) {
	if m.offset == nil {
		return 0, false
	}
	return *m.offset, true
}

// GetColumnCastAndTypedSlice is a mock implementation that returns the cast
// type and formatted slice.
func (m mockQueryState) GetColumnCastAndTypedSlice(
	colName string,
	slice []any,
) (string, any, bool) {
	if colName == "id" {
		res := make([]int64, len(slice))
		for i, v := range slice {
			if v == nil {
				continue
			}
			res[i] = v.(int64)
		}
		return "::bigint[]", res, false
	}
	res := make([]string, len(slice))
	for i, v := range slice {
		if v == nil {
			continue
		}
		res[i] = v.(string)
	}
	return "::text[]", res, false
}
