package repository

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDB holds the database connection and container for integration tests
type TestDB struct {
	DB        *sql.DB
	Container testcontainers.Container
}

// SetupTestDatabase creates a PostgreSQL testcontainer and runs migrations
func SetupTestDatabase(ctx context.Context) (*TestDB, error) {
	dbName := "testdb"
	dbUser := "testuser"
	dbPassword := "testpass"

	container, err := tcpostgres.Run(ctx,
		"postgres:15-alpine",
		tcpostgres.WithDatabase(dbName),
		tcpostgres.WithUsername(dbUser),
		tcpostgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		db.Close()
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &TestDB{
		DB:        db,
		Container: container,
	}, nil
}

// Teardown cleans up the test database
func (t *TestDB) Teardown(ctx context.Context) {
	if t.DB != nil {
		t.DB.Close()
	}
	if t.Container != nil {
		t.Container.Terminate(ctx)
	}
}

// runMigrations applies all migrations to the database using goose
func runMigrations(db *sql.DB) error {
	// Get migrations path relative to this file
	_, currentFile, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations")

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, migrationsPath); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// TruncateTables clears all data from tables (useful between tests)
func (t *TestDB) TruncateTables(ctx context.Context) error {
	tables := []string{"appointments", "doctors", "services", "patients"}
	for _, table := range tables {
		if _, err := t.DB.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)); err != nil {
			return fmt.Errorf("failed to truncate %s: %w", table, err)
		}
	}
	// Reset sequences
	for _, table := range tables {
		if _, err := t.DB.ExecContext(ctx, fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1", table)); err != nil {
			return fmt.Errorf("failed to reset sequence for %s: %w", table, err)
		}
	}
	return nil
}
