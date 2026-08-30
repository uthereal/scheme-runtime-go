package contract

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DB defines the interface for database execution operations.
type DB interface {
	// Exec executes a query without returning any rows, typically an INSERT,
	// UPDATE, or DELETE.
	Exec(
		ctxExec context.Context,
		sql string,
		arguments ...any,
	) (pgconn.CommandTag, error)

	// Query executes a query that returns multiple rows, typically a SELECT.
	Query(
		ctxQuery context.Context,
		sql string,
		arguments ...any,
	) (pgx.Rows, error)

	// QueryRow executes a query that is expected to return at most one row.
	QueryRow(
		ctxQueryRow context.Context,
		sql string,
		arguments ...any,
	) pgx.Row
}
