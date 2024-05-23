package chat

import "github.com/kavkaco/Kavka-Core/internal/model"

type GetChatValidation struct {
	ChatID model.ChatID `validate:"required"`
}

type GetUserChatsValidation struct {
	UserID model.UserID `validate:"required"`
}

type CreateDirectValidation struct {
	UserID          model.UserID `validate:"required"`
	RecipientUserID model.UserID `validate:"required"`
}

type CreateChannelValidation struct {
	UserID      model.UserID `validate:"required"`
	Title       string       `validate:"required"`
	Username    string       `validate:"required"`
	Description string       `validate:"required"`
}
type CreateGroupValidation struct {
	UserID      model.UserID `validate:"required"`
	Title       string       `validate:"required"`
	Username    string       `validate:"required"`
	Description string       `validate:"required"`
}
