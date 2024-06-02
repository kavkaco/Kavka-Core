package repository

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyTaken = errors.New("email already taken")
	ErrChatNotFound      = errors.New("chat not found")
	ErrChatAlreadyExists = errors.New("chat already exists")
	ErrMessageNotFound   = errors.New("message not found")
	ErrAuthNotFound      = errors.New("auth not found")
)

var (
	ErrNotModified = errors.New("document not modified")
	ErrNotDeleted  = errors.New("document not deleted")
)
