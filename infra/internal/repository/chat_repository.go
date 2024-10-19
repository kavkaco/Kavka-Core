package repository

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
)

type ChatRepository interface {
	GetChat(ctx context.Context, chatID model.ChatID) (*model.Chat, error)
	Create(ctx context.Context, chatModel model.Chat) (*model.Chat, error)
	Destroy(ctx context.Context, chatID model.ChatID) error
	GetUserChats(ctx context.Context, userID model.UserID, chatIDs []model.ChatID) ([]model.ChatDTO, error)
	GetDirectChat(ctx context.Context, userID model.UserID, recipientUserID model.UserID) (*model.Chat, error)
	GetChatMembers(chatID model.ChatID) []model.Member
	JoinChat(ctx context.Context, chatType string, userID string, chatID model.ChatID) error
	AddToUsersChatsList(ctx context.Context, userID string, chatID model.ChatID) error
}
