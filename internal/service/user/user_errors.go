package message

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUpdateUser   = errors.New("failed to update user")
)
