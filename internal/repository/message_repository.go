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
	FetchMessagesPaginated(ctx context.Context, chatID model.ChatID, skip, limit int64) ([]*model.MessageGetter, error)
	UpdateMessageContent(ctx context.Context, chatID model.ChatID, messageID model.MessageID, newMessageContent string) error
	Delete(ctx context.Context, chatID model.ChatID, messageID model.MessageID) error

	InsertV2(ctx context.Context, doc *model.MessageDocumentV2) (*model.MessageDocumentV2, error)
	FetchMessagesV2(ctx context.Context, chatID model.ChatID, skip, limit int64) ([]*model.MessageDocumentV2, error)
	FetchMessageV2(ctx context.Context, chatID model.ChatID, messageID model.MessageID) (*model.MessageDocumentV2, error)
	FetchLastMessageV2(ctx context.Context, chatID model.ChatID) (*model.MessageDocumentV2, error)
	UpdateMessageContentV2(ctx context.Context, chatID model.ChatID, messageID model.MessageID, newMessageContent string) error
	DeleteV2(ctx context.Context, chatID model.ChatID, messageID model.MessageID) error
	CountMessagesV2(ctx context.Context, chatID model.ChatID) (int64, error)
}
