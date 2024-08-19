package message

import "github.com/kavkaco/Kavka-Core/internal/model"

type InsertTextMessageValidation struct {
	ChatID         model.ChatID `validate:"required"`
	UserID         model.UserID `validate:"required"`
	MessageContent string       `validate:"required"`
}

type DeleteMessageValidation struct {
	ChatID    model.ChatID    `validate:"required"`
	UserID    model.UserID    `validate:"required"`
	MessageID model.MessageID `validate:"required"`
}
