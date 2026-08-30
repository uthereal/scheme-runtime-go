package grammar

import (
	"strings"

	"github.com/jackc/pgx/v5"
)

// PostgresGrammar compiles QueryStateProvider nodes into raw
// Postgres-compatible SQL strings and bindings.
type PostgresGrammar struct{}

// NewPostgresGrammar constructs a new PostgresGrammar.
func NewPostgresGrammar() *PostgresGrammar {
	return &PostgresGrammar{}
}

// sanitizeColumn sanitizes an identifier string for use in PostgreSQL query
// construction by escaping parts and splitting on dots.
func sanitizeColumn(identifier string) string {
	if identifier == "*" {
		return "*"
	}

	before, ok := strings.CutSuffix(identifier, ".*")
	if ok {
		if !strings.Contains(before, ".") {
			return pgx.Identifier{before}.Sanitize() + ".*"
		}
		partsBefore := strings.Split(before, ".")
		return pgx.Identifier(partsBefore).Sanitize() + ".*"
	}

	if !strings.Contains(identifier, ".") {
		return pgx.Identifier{identifier}.Sanitize()
	}

	partsIdentifier := strings.Split(identifier, ".")
	return pgx.Identifier(partsIdentifier).Sanitize()
}
