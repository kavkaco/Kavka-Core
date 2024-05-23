package repository

import "errors"

var (
	ErrChatNotFound      = errors.New("chat not found")
	ErrChatAlreadyExists = errors.New("chat already exists")
	ErrMessageNotFound   = errors.New("message not found")
)

var (
	ErrNotModified = errors.New("document not modified")
	ErrNotDeleted  = errors.New("document not deleted")
)
