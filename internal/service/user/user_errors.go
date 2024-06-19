package message

import "errors"

var (
	ErrNotFound   = errors.New("user not found")
	ErrUpdateUser = errors.New("failed to update user")
)
