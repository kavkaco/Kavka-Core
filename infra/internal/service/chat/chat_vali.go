package chat

import "github.com/kavkaco/Kavka-Core/internal/model"

type getChatValidation struct {
	ChatID model.ChatID `validate:"required"`
}

type getUserChatsValidation struct {
	UserID model.UserID `validate:"required"`
}

type createDirectValidation struct {
	UserID          model.UserID `validate:"required"`
	RecipientUserID model.UserID `validate:"required"`
}

type createChannelValidation struct {
	UserID      model.UserID `validate:"required"`
	Title       string       `validate:"required,min=1"`
	Username    string       `validate:"required,min=3"`
	Description string
}
type createGroupValidation struct {
	UserID      model.UserID `validate:"required"`
	Title       string       `validate:"required,min=1"`
	Username    string       `validate:"required,min=3"`
	Description string
}
