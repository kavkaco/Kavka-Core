package user

import (
	"errors"
	"os/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrEmailAlreadyTaken    = errors.New("email already taken")
)

type UserRepository interface {
	FindByID(staticID primitive.ObjectID) (*user.User, error)
	Where(filter any) ([]*user.User, error)
	Create(u *CreateUserData) (*user.User, error)
}
