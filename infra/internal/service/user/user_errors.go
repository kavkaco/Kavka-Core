package user

import "errors"

var (
	ErrNotFound   = errors.New("user not found")
	ErrUpdateUser = errors.New("failed to update user")
	ErrDeleteUser = errors.New("failed to delete user")
)
