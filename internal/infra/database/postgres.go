package database

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"zpmeow/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// ConnectDB connects to the PostgreSQL database and returns a connection instance
func ConnectDB(cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSslMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// RunMigrations executes the database migrations
func RunMigrations(db *sqlx.DB) error {
	migrationsDir := "internal/database/migrations"

	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("could not read migrations directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			filePath := filepath.Join(migrationsDir, file.Name())
			query, err := ioutil.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("could not read migration file %s: %w", file.Name(), err)
			}

			if _, err := db.Exec(string(query)); err != nil {
				return fmt.Errorf("failed to execute migration file %s: %w", file.Name(), err)
			}
		}
	}

	return nil
}
