package migration

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Migrate(migrationsPath string) error {
	// Create schema_migrations table if it doesn't exist
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version bigint PRIMARY KEY,
			dirty boolean NOT NULL DEFAULT FALSE
		);
	`
	if _, err := m.db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Get the current migration version
	var currentVersion int64
	err := m.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&currentVersion)
	if err != nil {
		return fmt.Errorf("failed to query current migration version: %w", err)
	}

	// Read all migration files
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []struct {
		version int64
		file    string
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".up.sql") {
			continue
		}
		versionStr := strings.Split(file.Name(), "_")[0]
		version, err := strconv.ParseInt(versionStr, 10, 64)
		if err != nil {
			continue
		}
		migrations = append(migrations, struct {
			version int64
			file    string
		}{version, file.Name()})
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].version < migrations[j].version
	})

	// Execute pending migrations
	for _, mig := range migrations {
		if mig.version <= currentVersion {
			continue
		}

		sqlFile := filepath.Join(migrationsPath, mig.file)
		sqlBytes, err := ioutil.ReadFile(sqlFile)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", mig.file, err)
		}

		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", mig.file, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version, dirty) VALUES ($1, false)", mig.version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to track migration version: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration transaction: %w", err)
		}

		fmt.Printf("Migration %d executed successfully\n", mig.version)
	}

	return nil
}
