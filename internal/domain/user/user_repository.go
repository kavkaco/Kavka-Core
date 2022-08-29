package user

import (
	"errors"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrEmailAlreadyTaken    = errors.New("email already taken")
)

type Repository interface {
	Create(name string, lastName string, username string, passwordHash string, email string) (*User, error)
	GetByStaticUUID(staticID string) (*User, error)
	Delete(staticID string) error
	// TODO - complete: edit profile
}
