package repository

import (
	"context"
	"time"

	"github.com/kavkaco/Kavka-Core/internal/model"
)

type AuthRepository interface {
	Create(ctx context.Context, authModel *model.Auth) (*model.Auth, error)
	GetUserAuth(ctx context.Context, userID model.UserID) (*model.Auth, error)
	ChangePassword(ctx context.Context, userID model.UserID, passwordHash string) error
	VerifyEmail(ctx context.Context, userID model.UserID) error
	IncrementFailedLoginAttempts(ctx context.Context, userID model.UserID) error
	ClearFailedLoginAttempts(ctx context.Context, userID model.UserID) error
	LockAccount(ctx context.Context, userID model.UserID, lockDuration time.Duration) error
	UnlockAccount(ctx context.Context, userID model.UserID) error
}

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

type MessageRepository interface {
	FindMessage(ctx context.Context, chatID model.ChatID, messageID model.MessageID) (*model.Message, error)
	Create(ctx context.Context, chatID model.ChatID) error
	FetchMessages(ctx context.Context, chatID model.ChatID) ([]model.Message, error)
	Insert(ctx context.Context, chatID model.ChatID, message *model.Message) (*model.Message, error)
	UpdateMessageContent(ctx context.Context, chatID model.ChatID, messageID model.MessageID, newMessageContent string) error
	Delete(ctx context.Context, chatID model.ChatID, messageID model.MessageID) error
}

type UserRepository interface {
	GetChats(ctx context.Context, userID model.UserID) ([]model.ChatID, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	AddToUserChats(ctx context.Context, userID model.UserID, chatID model.ChatID) error
	Update(ctx context.Context, userID string, name, lastName, username, biography string) error
	FindByUserID(ctx context.Context, userID model.UserID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}
