package user

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrEmailAlreadyTaken    = errors.New("email already taken")
)

type Repository interface {
	Create(name string, lastName string, username string, passwordHash string, email string) (*User, error)
	GetByStaticUUID(staticID uuid.UUID) (*User, error)
	Delete(staticID uuid.UUID) error
	// TODO - complete: edit profile
}
