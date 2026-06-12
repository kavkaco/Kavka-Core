package repo_sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kavkaco/Kavka-Core/database"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type SQLAdapter struct {
	DB     *sql.DB
	Driver database.DriverType
}

func NewSQLAdapter(cfg *database.DatabaseConfig) (*SQLAdapter, error) {
	var db *sql.DB
	var err error

	switch cfg.Driver {
	case database.DriverPostgres:
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return nil, fmt.Errorf("postgres connection failed: %w", err)
		}

	case database.DriverSQLite:
		db, err = sql.Open("sqlite3", cfg.SQLitePath)
		if err != nil {
			return nil, fmt.Errorf("sqlite connection failed: %w", err)
		}

	default:
		return nil, database.ErrUnsupportedDriver
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	adapter := &SQLAdapter{
		DB:     db,
		Driver: cfg.Driver,
	}

	if err = adapter.runMigrations(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return adapter, nil
}

func (a *SQLAdapter) runMigrations() error {
	switch a.Driver {
	case database.DriverPostgres:
		_, err := a.DB.Exec(SchemaPostgres)
		return err
	case database.DriverSQLite:
		_, err := a.DB.Exec(SchemaSQLite)
		return err
	}
	return nil
}

func (a *SQLAdapter) Close() error {
	return a.DB.Close()
}
