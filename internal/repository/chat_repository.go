package repository

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatRepository interface {
	GetChat(ctx context.Context, chatID model.ChatID) (*model.Chat, error)
	Create(ctx context.Context, chatModel model.Chat) (*model.Chat, error)
	Destroy(ctx context.Context, chatID model.ChatID) error
	GetUserChats(ctx context.Context, chatIDs []model.ChatID) ([]model.ChatDTO, error)
	FindBySides(ctx context.Context, userID model.UserID, recipientUserID model.UserID) (*model.Chat, error)
	GetChatMembers(chatID model.ChatID) []model.Member
	JoinChat(ctx context.Context, userID string, chatID primitive.ObjectID) error
}
