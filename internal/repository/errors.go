package repository

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrNotModified      = errors.New("document not modified")
	ErrNotDeleted       = errors.New("document not deleted")
	ErrUniqueConstraint = errors.New("unique constraint violation")
)
