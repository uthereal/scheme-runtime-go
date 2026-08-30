package contract

// SQLCompiler defines the behavior for translating structured query builders
// and mutation parameters into specific SQL statements and data bindings.
type SQLCompiler interface {
	// CompileSelect translates a read-only query representation into a SQL
	// SELECT statement and its parameter bindings.
	CompileSelect(
		state QueryStateProvider,
	) (sql string, bindings []any)

	// CompileInsert translates a series of column value records into a SQL
	// INSERT statement.
	CompileInsert(
		state QueryStateProvider,
		values [][]ColumnValue,
	) (sql string, bindings []any)

	// CompileInsertReturning compiles a SQL INSERT statement that returns a
	// specific subset of projected column fields.
	CompileInsertReturning(
		state QueryStateProvider,
		values [][]ColumnValue,
		returning []string,
	) (sql string, bindings []any)

	// CompileUpdate translates a set of column changes into a SQL UPDATE
	// statement filtered by the query constraints.
	CompileUpdate(
		state QueryStateProvider,
		values []ColumnValue,
	) (sql string, bindings []any)

	// CompileUpdateReturning compiles a SQL UPDATE statement that returns a
	// specific subset of projected column fields.
	CompileUpdateReturning(
		state QueryStateProvider,
		values []ColumnValue,
		returning []string,
	) (sql string, bindings []any)

	// CompileDelete compiles a SQL DELETE statement filtered by the query
	// state constraints.
	CompileDelete(
		state QueryStateProvider,
	) (sql string, bindings []any)

	// CompileDeleteReturning compiles a SQL DELETE statement that returns a
	// specific subset of projected column fields.
	CompileDeleteReturning(
		state QueryStateProvider,
		returning []string,
	) (sql string, bindings []any)

	// CompileUpsert translates records into an INSERT statement with an ON
	// CONFLICT clause to handle insert collisions as updates.
	CompileUpsert(
		state QueryStateProvider,
		values [][]ColumnValue,
		conflictColumns []string,
	) (sql string, bindings []any)
}
