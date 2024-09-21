package message

import "github.com/kavkaco/Kavka-Core/internal/model"

type insertTextMessageValidation struct {
	ChatID         model.ChatID `validate:"required"`
	UserID         model.UserID `validate:"required"`
	MessageContent string       `validate:"required"`
}

type deleteMessageValidation struct {
	ChatID    model.ChatID    `validate:"required"`
	UserID    model.UserID    `validate:"required"`
	MessageID model.MessageID `validate:"required"`
}
