package db

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// MigrationManager handles database migrations
type MigrationManager struct {
	db         *sql.DB
	migrations string // Path to migrations directory
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *sql.DB, migrationsPath string) *MigrationManager {
	return &MigrationManager{
		db:         db,
		migrations: migrationsPath,
	}
}

// ensureMigrationTable creates the migrations tracking table if it doesn't exist
func (m *MigrationManager) ensureMigrationTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`
	_, err := m.db.ExecContext(ctx, query)
	return err
}

// getAppliedMigrations returns a map of applied migration versions
func (m *MigrationManager) getAppliedMigrations(ctx context.Context) (map[string]bool, error) {
	applied := make(map[string]bool)
	
	rows, err := m.db.QueryContext(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	return applied, rows.Err()
}

// RunMigrations executes all pending migrations
func (m *MigrationManager) RunMigrations(ctx context.Context) error {
	// Ensure migration table exists
	if err := m.ensureMigrationTable(ctx); err != nil {
		return fmt.Errorf("ensuring migration table: %w", err)
	}

	// Get applied migrations
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("getting applied migrations: %w", err)
	}

	// Read migration files
	files, err := ioutil.ReadDir(m.migrations)
	if err != nil {
		return fmt.Errorf("reading migrations directory: %w", err)
	}

	// Filter and sort migration files
	var migrations []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			migrations = append(migrations, f.Name())
		}
	}
	sort.Strings(migrations)

	// Start a transaction for all migrations
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting migration transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute each pending migration
	for _, migration := range migrations {
		version := strings.TrimSuffix(migration, ".sql")
		
		// Skip if already applied
		if applied[version] {
			continue
		}

		// Read migration file
		content, err := ioutil.ReadFile(filepath.Join(m.migrations, migration))
		if err != nil {
			return fmt.Errorf("reading migration %s: %w", migration, err)
		}

		// Execute migration
		if _, err := tx.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("executing migration %s: %w", migration, err)
		}

		// Record migration
		if _, err := tx.ExecContext(ctx, 
			"INSERT INTO schema_migrations (version, applied_at) VALUES ($1, $2)",
			version, time.Now()); err != nil {
			return fmt.Errorf("recording migration %s: %w", migration, err)
		}
	}

	// Commit all migrations
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing migrations: %w", err)
	}

	return nil
}