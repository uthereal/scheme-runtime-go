package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_PostgresGrammar_NewPostgresGrammar tests NewPostgresGrammar constructor.
func Test_PostgresGrammar_NewPostgresGrammar(t *testing.T) {
	g := NewPostgresGrammar()
	assert.NotNil(t, g)
}

// Test_PostgresGrammar_SanitizeColumn tests sanitizeColumn identifier
// sanitation.
func Test_PostgresGrammar_SanitizeColumn(t *testing.T) {
	// 1. Star column
	assert.Equal(t, "*", sanitizeColumn("*"))

	// 2. Column suffix ".*" without schema prefix
	assert.Equal(t, `"users".*`, sanitizeColumn("users.*"))

	// 3. Column suffix ".*" with schema prefix
	assert.Equal(t, `"public"."users".*`, sanitizeColumn("public.users.*"))

	// 4. Identifier without dot
	assert.Equal(t, `"id"`, sanitizeColumn("id"))

	// 5. Identifier with dot
	assert.Equal(t, `"users"."id"`, sanitizeColumn("users.id"))
}
