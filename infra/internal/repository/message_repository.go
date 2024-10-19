package repository

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
)

type MessageRepository interface {
	Create(ctx context.Context, chatID model.ChatID) error
	Insert(ctx context.Context, chatID model.ChatID, message *model.Message) (*model.Message, error)
	FetchLastMessage(ctx context.Context, chatID model.ChatID) (*model.Message, error)
	FetchMessage(ctx context.Context, chatID model.ChatID, messageID model.MessageID) (*model.Message, error)
	FetchMessages(ctx context.Context, chatID model.ChatID) ([]*model.MessageGetter, error)
	UpdateMessageContent(ctx context.Context, chatID model.ChatID, messageID model.MessageID, newMessageContent string) error
	Delete(ctx context.Context, chatID model.ChatID, messageID model.MessageID) error
}
