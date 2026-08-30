package orm

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/column"
)

// mockDb is a mock database implementation for ORM testing.
type mockDb struct {
	// execResult is the result returned by Exec.
	execResult pgconn.CommandTag
	// execErr is the error returned by Exec.
	execErr error
	// queryRows is the mock rows returned by Query.
	queryRows *mockRows
	// queryErr is the error returned by Query.
	queryErr error
	// queryRowValue is the mock row returned by QueryRow.
	queryRowValue *mockRow
}

// mockRow is a mock row for database scanning.
type mockRow struct {
	// value is the single mocked int64 value.
	value int64
	// err is the mocked scan error.
	err error
}

// mockRows is a mock of pgx.Rows for testing database queries.
type mockRows struct {
	// cols is the names of the mock columns.
	cols []string
	// records contains the mocked query result rows.
	records [][]any
	// cursor is the current row index.
	cursor int
	// closed indicates if the rows were closed.
	closed bool
	// rowsErr is the error returned by Err.
	rowsErr error
}

// testModel is a simple model structure for testing ORM capabilities.
type testModel struct {
	// ID is the unique identifier of the test model.
	ID int64 `db:"id"`
	// Email is the email address of the test model.
	Email string `db:"email"`
}

// testMutator is a mutator structure for the test model.
type testMutator struct {
	// ID is the unique identifier pointer.
	ID *int64 `db:"id"`
	// Email is the email address pointer.
	Email *string `db:"email"`
}

// testSchema defines columns for the test model.
var testSchema = struct {
	ID    column.NumericColumn[testModel, int64]
	Email column.StringColumn[testModel, string]
	Age   column.NullableNumericColumn[testModel, int]
}{
	ID: column.NumericColumn[testModel, int64]{
		Name: "id",
	},
	Email: column.StringColumn[testModel, string]{
		Name: "email",
	},
	Age: column.NullableNumericColumn[testModel, int]{
		Name: "age",
	},
}

// testTable specifies the database table metadata for testModel.
var testTable = contract.TableMetadata[testModel]{
	SchemaName: "public",
	TableName:  "users",
	DefaultColumns: []contract.Column[testModel]{
		testSchema.ID,
		testSchema.Email,
	},
}

// Exec executes a SQL command against the mock database.
func (m *mockDb) Exec(
	_ context.Context,
	_ string,
	_ ...any,
) (pgconn.CommandTag, error) {
	return m.execResult, m.execErr
}

// Query runs a query and returns mock rows.
func (m *mockDb) Query(
	_ context.Context,
	_ string,
	_ ...any,
) (pgx.Rows, error) {
	if m.queryErr != nil {
		return nil, m.queryErr
	}
	return m.queryRows, nil
}

// QueryRow runs a query expecting a single row return.
func (m *mockDb) QueryRow(
	_ context.Context,
	_ string,
	_ ...any,
) pgx.Row {
	return m.queryRowValue
}

// Scan scans mock database columns into destination variables.
func (r *mockRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) == 0 {
		return nil
	}
	pInt, okInt := dest[0].(*int64)
	if okInt {
		*pInt = r.value
		return nil
	}
	pUint, okUint := dest[0].(*uint64)
	if okUint {
		*pUint = uint64(r.value)
		return nil
	}
	pFloat, okFloat := dest[0].(*float64)
	if okFloat {
		*pFloat = float64(r.value)
		return nil
	}
	pAny, okAny := dest[0].(*any)
	if okAny {
		*pAny = r.value
		return nil
	}
	return nil
}

// Close closes the mock rows cursor.
func (m *mockRows) Close() {
	m.closed = true
}

// CommandTag returns the command tag for the query.
func (m *mockRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

// Conn returns the underlying connection.
func (m *mockRows) Conn() *pgx.Conn {
	return nil
}

// RawValues returns raw values for the row.
func (m *mockRows) RawValues() [][]byte {
	return nil
}

// Err returns any query error.
func (m *mockRows) Err() error {
	return m.rowsErr
}

// Next advances the cursor to the next record.
func (m *mockRows) Next() bool {
	if m.cursor < len(m.records) {
		m.cursor++
		return true
	}
	return false
}

// Scan scans columns of the current row into destination values.
func (m *mockRows) Scan(dest ...any) error {
	if m.cursor-1 >= len(m.records) {
		return errors.New("out of bounds scan")
	}
	row := m.records[m.cursor-1]
	for i, val := range row {
		if i >= len(dest) {
			break
		}
		reflectValue := reflectHydrate(dest[i], val)
		if !reflectValue {
			return errors.New("failed scanning type")
		}
	}
	return nil
}

// Values returns the values of the current row.
func (m *mockRows) Values() ([]any, error) {
	return nil, nil
}

// FieldDescriptions returns mock field descriptions.
func (m *mockRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

// reflectHydrate populates destination pointer with source value.
func reflectHydrate(
	dest any,
	src any,
) bool {
	switch d := dest.(type) {
	case *int64:
		s, ok := src.(int64)
		if ok {
			*d = s
			return true
		}
	case *string:
		s, ok := src.(string)
		if ok {
			*d = s
			return true
		}
	case *any:
		*d = src
		return true
	}
	return false
}

// testModelHydrate hydrates a testModel with the specified columns.
func testModelHydrate(
	model *testModel,
	columns []string,
) []any {
	pointers := make([]any, len(columns))
	for i, col := range columns {
		switch col {
		case "id":
			pointers[i] = &model.ID
		case "email":
			pointers[i] = &model.Email
		default:
			var discard any
			pointers[i] = &discard
		}
	}
	return pointers
}

// testMutatorDehydrate dehydrates a testMutator into columns and values.
func testMutatorDehydrate(
	mutator *testMutator,
) ([]string, []any) {
	var cols []string
	var vals []any
	if mutator.ID != nil {
		cols = append(cols, "id")
		vals = append(vals, *mutator.ID)
	}
	if mutator.Email != nil {
		cols = append(cols, "email")
		vals = append(vals, *mutator.Email)
	}
	return cols, vals
}
