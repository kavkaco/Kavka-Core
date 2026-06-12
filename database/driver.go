package database

import "errors"

type DriverType string

const (
	DriverMongo    DriverType = "mongodb"
	DriverPostgres DriverType = "postgres"
	DriverSQLite   DriverType = "sqlite"
)

var (
	ErrUnsupportedDriver = errors.New("unsupported database driver")
	ErrDriverNotInitialized = errors.New("database driver not initialized")
)

type DatabaseConfig struct {
	Driver   DriverType `koanf:"driver"`
	Host     string     `koanf:"host"`
	Port     int        `koanf:"port"`
	Username string     `koanf:"username"`
	Password string     `koanf:"password"`
	DBName   string     `koanf:"db_name"`
	SQLitePath string   `koanf:"sqlite_path"`
	SSLMode  string     `koanf:"ssl_mode"`
}
