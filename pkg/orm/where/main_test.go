package where

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// mockQueryStateProvider is a mock implementation of
// contract.QueryStateProvider for testing.
type mockQueryStateProvider struct{}

// GetSchemaName returns the mock schema name.
func (m mockQueryStateProvider) GetSchemaName() string {
	return "public"
}

// GetTableName returns the mock physical table name.
func (m mockQueryStateProvider) GetTableName() string {
	return "users"
}

// GetDefaultColumns returns the mock standard select columns.
func (m mockQueryStateProvider) GetDefaultColumns() []string {
	return []string{"id"}
}

// GetSelectedColumns returns the mock columns to be fetched.
func (m mockQueryStateProvider) GetSelectedColumns() []string {
	return []string{"id"}
}

// IsDistinct returns true if mock SELECT DISTINCT is requested.
func (m mockQueryStateProvider) IsDistinct() bool {
	return false
}

// GetAggregate returns the mock select aggregate function state.
func (m mockQueryStateProvider) GetAggregate() *contract.AggregateState {
	return nil
}

// GetWheres returns mock applied query WHERE conditions.
func (m mockQueryStateProvider) GetWheres() []contract.Where {
	return nil
}

// GetOrders returns mock applied sorting ORDER BY clauses.
func (m mockQueryStateProvider) GetOrders() []contract.Order {
	return nil
}

// GetGroups returns the mock list of columns to GROUP BY.
func (m mockQueryStateProvider) GetGroups() []string {
	return nil
}

// GetHavings returns mock applied grouped HAVING conditions.
func (m mockQueryStateProvider) GetHavings() []contract.Where {
	return nil
}

// GetLimit returns mock maximum number of rows to return.
func (m mockQueryStateProvider) GetLimit() (uint64, bool) {
	return 0, false
}

// GetOffset returns the mock number of rows to skip.
func (m mockQueryStateProvider) GetOffset() (uint64, bool) {
	return 0, false
}

// GetColumnCastAndTypedSlice returns mock cast and slice.
func (m mockQueryStateProvider) GetColumnCastAndTypedSlice(
	_ string,
	slice []any,
) (string, any, bool) {
	return "::text[]", slice, false
}
