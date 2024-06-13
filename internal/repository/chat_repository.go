package repository

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
)

type ChatRepository interface {
	SearchInChats(ctx context.Context, key string) ([]model.Chat, error)
	UpdateChatLastMessage(ctx context.Context, chatID model.ChatID, lastMessage model.LastMessage) error
	Create(ctx context.Context, chatModel model.Chat) (*model.Chat, error)
	Destroy(ctx context.Context, chatID model.ChatID) error
	FindManyByChatID(ctx context.Context, chatIDs []model.ChatID) ([]model.Chat, error)
	FindByID(ctx context.Context, chatID model.ChatID) (*model.Chat, error)
	FindBySides(ctx context.Context, sides [2]model.UserID) (*model.Chat, error)
	GetChatMembers(chatID model.ChatID) []model.Member
}
