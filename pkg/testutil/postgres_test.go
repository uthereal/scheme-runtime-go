package testutil_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uthereal/scheme-runtime-go/pkg/testutil"
)

// TestPostgresContainer_Lifecycle verifies that a PostgreSQL container can be
// started, configured with a template DB, cloned, and stopped successfully.
func TestPostgresContainer_Lifecycle(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(
		context.Background(),
		2*time.Minute,
	)
	defer cancel()

	pgContainer, err := testutil.StartPostgresContainer(ctxTimeout)
	require.NoError(t, err)

	t.Cleanup(func() {
		ctxStop, cancelStop := context.WithTimeout(
			context.Background(),
			30*time.Second,
		)
		defer cancelStop()
		errStop := testutil.StopPostgresContainer(ctxStop, pgContainer)
		assert.NoError(t, errStop)
	})

	assert.NotNil(t, pgContainer)

	schemaDDL := `CREATE TABLE test_table (id INT PRIMARY KEY);`
	err = pgContainer.SetupTemplateDB(
		ctxTimeout,
		"test_template_db",
		schemaDDL,
	)
	require.NoError(t, err)

	db, cleanup, err := pgContainer.CreateDB(ctxTimeout, t)
	require.NoError(t, err)
	require.NotNil(t, db)

	t.Cleanup(func() {
		cleanup()
	})

	_, err = db.Exec(
		ctxTimeout,
		"INSERT INTO test_table (id) VALUES (1)",
	)
	assert.NoError(t, err)
}
