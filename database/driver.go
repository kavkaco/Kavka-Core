package database

import "errors"

var (
	ErrUnsupportedDriver    = errors.New("unsupported database driver")
	ErrDriverNotInitialized = errors.New("database driver not initialized")
)
