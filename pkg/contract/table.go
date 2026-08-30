package contract

// TableMetadata provides static metadata about a database table.
type TableMetadata[Model any] struct {
	// SchemaName is the name of the schema where the table resides.
	SchemaName string

	// TableName is the name of the database table.
	TableName string

	// DefaultColumns is the slice of default columns to select from the table.
	DefaultColumns []Column[Model]
}
