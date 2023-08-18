package user

import (
	"errors"
	"os/user"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrEmailAlreadyTaken    = errors.New("email already taken")
)

type Repository interface {
	Create(phone string) (*User, error)
	GetByStaticUUID(staticID string) (*User, error)
	Delete(staticID string) error
	Update(u *user.User) (*User, error)
}
