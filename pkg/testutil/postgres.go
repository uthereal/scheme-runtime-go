package testutil

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// PostgresContainer represents a PostgreSQL container for testing.
type PostgresContainer struct {
	*tcpostgres.PostgresContainer

	mu             sync.RWMutex
	dbNameTemplate *string
}

// postgresImage defines the Docker image for PostgreSQL.
const postgresImage = "postgres:18-alpine"

// StartPostgresContainer starts a new PostgreSQL container using
// testcontainers.
func StartPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := tcpostgres.Run(
		ctx,
		postgresImage,
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("user"),
		tcpostgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog(
				"database system is ready to accept connections",
			).WithOccurrence(2).WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to start postgres container -> %w",
			err,
		)
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
	}, nil
}

// StopPostgresContainer stops a PostgreSQL container.
func StopPostgresContainer(
	ctxStop context.Context,
	pgContainer *PostgresContainer,
) error {
	return pgContainer.Terminate(ctxStop)
}

// SetupTemplateDB creates a template database with the initial schema DDL.
func (c *PostgresContainer) SetupTemplateDB(
	ctx context.Context,
	dbNameTemplate string,
	schemaDDL string,
) error {
	ctxAdmin, cancelAdmin := context.WithTimeout(ctx, 10*time.Second)
	defer cancelAdmin()

	rawAdminConn, err := c.ConnectionString(ctxAdmin, "sslmode=disable")
	if err != nil {
		return fmt.Errorf(
			"failed to get base admin connection string -> %w",
			err,
		)
	}

	adminURL, err := url.Parse(rawAdminConn)
	if err != nil {
		return fmt.Errorf("failed to parse connection string -> %w", err)
	}

	adminURL.Path = "/postgres"
	adminConn := adminURL.String()

	adminDB, err := sql.Open("pgx", adminConn)
	if err != nil {
		return fmt.Errorf("failed to open admin database -> %w", err)
	}
	defer func(adminDB *sql.DB) {
		_ = adminDB.Close()
	}(adminDB)

	safeTemplateDBName := pgx.Identifier{dbNameTemplate}.Sanitize()
	createQ := fmt.Sprintf("CREATE DATABASE %s", safeTemplateDBName)
	_, err = adminDB.ExecContext(ctxAdmin, createQ)
	if err != nil {
		return fmt.Errorf(
			"failed to create template db %s -> %w",
			dbNameTemplate,
			err,
		)
	}

	ctxTarget, cancelTarget := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTarget()

	targetURL, err := url.Parse(rawAdminConn)
	if err != nil {
		return fmt.Errorf(
			"failed to parse connection string -> %w",
			err,
		)
	}

	targetURL.Path = "/" + dbNameTemplate
	targetConn := targetURL.String()

	targetDB, err := sql.Open("pgx", targetConn)
	if err != nil {
		return fmt.Errorf(
			"failed to open template database -> %w",
			err,
		)
	}
	defer func(targetDB *sql.DB) {
		_ = targetDB.Close()
	}(targetDB)

	_, err = targetDB.ExecContext(ctxTarget, schemaDDL)
	if err != nil {
		return fmt.Errorf(
			"failed to initialize template database schema -> %w",
			err,
		)
	}

	err = targetDB.Close()
	if err != nil {
		return fmt.Errorf(
			"failed to close template database connection -> %w",
			err,
		)
	}

	ctxTemplate, cancelTemplate := context.WithTimeout(ctx, 10*time.Second)
	defer cancelTemplate()

	alterQ := fmt.Sprintf(
		"ALTER DATABASE %s IS_TEMPLATE = true",
		safeTemplateDBName,
	)
	_, err = adminDB.ExecContext(ctxTemplate, alterQ)
	if err != nil {
		return fmt.Errorf(
			"failed to set template db %s to template -> %w",
			dbNameTemplate,
			err,
		)
	}

	c.mu.Lock()
	c.dbNameTemplate = &dbNameTemplate
	c.mu.Unlock()

	return nil
}

// CreateDB creates a new isolated database cloned from the template.
// It returns the pgx connection, cleanup function, and any error.
func (c *PostgresContainer) CreateDB(
	ctx context.Context,
	t *testing.T,
) (contract.DB, func(), error) {
	t.Helper()

	c.mu.RLock()
	templateName := c.dbNameTemplate
	c.mu.RUnlock()

	if templateName == nil {
		return nil,
			nil,
			errors.New(
				"template database not initialized -> call SetupTemplateDB first",
			)
	}

	hash := md5.Sum([]byte(t.Name()))
	dbName := fmt.Sprintf("test_db_%s", hex.EncodeToString(hash[:]))

	ctxAdmin, cancelAdmin := context.WithTimeout(ctx, 10*time.Second)
	defer cancelAdmin()

	rawAdminConn, err := c.ConnectionString(ctxAdmin, "sslmode=disable")
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to get base admin connection string -> %w",
			err,
		)
	}

	adminURL, err := url.Parse(rawAdminConn)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to parse connection string -> %w",
			err,
		)
	}

	adminURL.Path = "/postgres"
	adminConn := adminURL.String()

	adminDB, err := sql.Open("pgx", adminConn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open admin db -> %w", err)
	}
	defer func(adminDB *sql.DB) {
		_ = adminDB.Close()
	}(adminDB)

	safeDBName := pgx.Identifier{dbName}.Sanitize()
	safeTemplateName := pgx.Identifier{*templateName}.Sanitize()
	cmd := fmt.Sprintf(
		"CREATE DATABASE %s WITH TEMPLATE %s",
		safeDBName,
		safeTemplateName,
	)
	_, err = adminDB.ExecContext(ctxAdmin, cmd)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to clone template -> %w", err)
	}

	ctxTarget, cancelTarget := context.WithTimeout(ctx, 10*time.Second)
	defer cancelTarget()

	rawConn, err := c.ConnectionString(ctxTarget, "sslmode=disable")
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to get base conn string -> %w",
			err,
		)
	}

	u, err := url.Parse(rawConn)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to parse conn string URL -> %w",
			err,
		)
	}

	u.Path = "/" + dbName
	newConnStr := u.String()

	conn, err := pgx.Connect(ctxTarget, newConnStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to cloned db -> %w", err)
	}

	cleanup := func() {
		_ = conn.Close(context.Background())

		ctxCleanup, cancelCleanup := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancelCleanup()

		cleanupDB, err := sql.Open("pgx", adminConn)
		if err == nil {
			defer func(cleanupDB *sql.DB) {
				_ = cleanupDB.Close()
			}(cleanupDB)
			safeDropDBName := pgx.Identifier{dbName}.Sanitize()
			drop := fmt.Sprintf(
				"DROP DATABASE IF EXISTS %s WITH (FORCE)",
				safeDropDBName,
			)
			_, _ = cleanupDB.ExecContext(ctxCleanup, drop)
		}
	}

	return conn, cleanup, nil
}
