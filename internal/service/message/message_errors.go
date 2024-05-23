package message

import "errors"

var (
	ErrMessageNotFound = errors.New("message not found")
	ErrChatNotFound    = errors.New("chat not found")
	ErrDeleteMessage   = errors.New("failed to delete message")
	ErrInsertMessage   = errors.New("failed to insert message")
	ErrAccessDenied    = errors.New("access denied")
)
