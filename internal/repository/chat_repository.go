package repository

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatRepository interface {
	Create(ctx context.Context, chatModel model.Chat) (*model.Chat, error)
	Destroy(ctx context.Context, chatID model.ChatID) error
	GetUserChats(ctx context.Context, chatIDs []model.ChatID) ([]model.ChatGetter, error)
	FindByID(ctx context.Context, chatID model.ChatID) (*model.Chat, error)
	FindBySides(ctx context.Context, sides [2]model.UserID) (*model.Chat, error)
	GetChatMembers(chatID model.ChatID) []model.Member
	JoinChat(ctx context.Context, userID string, chatID primitive.ObjectID) error
}
