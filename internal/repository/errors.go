package repository

import "errors"

var (
	ErrChatNotFound      = errors.New("chat not found")
	ErrChatAlreadyExists = errors.New("chat already exists")
)

var ErrNotModified = errors.New("document not modified")
var ErrNotDeleted = errors.New("document not deleted")
