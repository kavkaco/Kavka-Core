package service

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model/chat"
	chatRepository "github.com/kavkaco/Kavka-Core/internal/repository/chat"
	userRepository "github.com/kavkaco/Kavka-Core/internal/repository/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type ChatService interface {
	GetChat(ctx context.Context, chatStaticID primitive.ObjectID) (*chat.ChatC, error)
	GetUserChats(ctx context.Context, userStaticID primitive.ObjectID) ([]chat.ChatC, error)
	CreateDirect(ctx context.Context, userStaticID primitive.ObjectID, targetStaticID primitive.ObjectID) (*chat.Chat, error)
	CreateGroup(ctx context.Context, userStaticID primitive.ObjectID, title string, username string, description string) (*chat.Chat, error)
	CreateChannel(ctx context.Context, userStaticID primitive.ObjectID, title string, username string, description string) (*chat.Chat, error)
}

type ChatManager struct {
	logger   *zap.Logger
	chatRepo chatRepository.ChatRepository
	userRepo userRepository.UserRepository
}

func NewChatService(logger *zap.Logger, chatRepo chatRepository.ChatRepository, userRepo userRepository.UserRepository) ChatService {
	return &ChatManager{logger, chatRepo, userRepo}
}

func (s *ChatManager) GetChat(ctx context.Context, chatStaticID primitive.ObjectID) (*chat.ChatC, error) {
	foundChat, err := s.chatRepo.FindChatOrSidesByStaticID(ctx, chatStaticID)
	if err != nil {
		return nil, err
	}

	return foundChat, nil
}

// FIXME - not implemented yet
func (s *ChatManager) GetUserChats(ctx context.Context, userStaticID primitive.ObjectID) ([]chat.ChatC, error) {
	return []chat.ChatC{}, nil
	// return s.chatRepo.GetChats(userStaticID)
}

func (s *ChatManager) CreateDirect(ctx context.Context, userStaticID primitive.ObjectID, targetStaticID primitive.ObjectID) (*chat.Chat, error) {
	sides := [2]primitive.ObjectID{
		userStaticID,
		targetStaticID,
	}

	// Check to do not be duplicated!
	dup, _ := s.chatRepo.FindBySides(ctx, sides)
	if dup != nil {
		return nil, chatRepository.ErrChatAlreadyExists
	}

	newChat := chat.NewChat(chat.TypeDirect, &chat.DirectChatDetail{
		Sides: sides,
	})

	return s.chatRepo.Create(ctx, *newChat)
}

func (s *ChatManager) CreateGroup(ctx context.Context, userStaticID primitive.ObjectID, title string, username string, description string) (*chat.Chat, error) {
	newChat := chat.NewChat(chat.TypeGroup, &chat.GroupChatDetail{
		Title:       title,
		Username:    username,
		Members:     []primitive.ObjectID{userStaticID},
		Admins:      []primitive.ObjectID{userStaticID},
		Description: description,
		Owner:       &userStaticID,
	})

	return s.chatRepo.Create(ctx, *newChat)
}

func (s *ChatManager) CreateChannel(ctx context.Context, userStaticID primitive.ObjectID, title string, username string, description string) (*chat.Chat, error) {
	newChat := chat.NewChat(chat.TypeGroup, &chat.GroupChatDetail{
		Title:       title,
		Username:    username,
		Members:     []primitive.ObjectID{userStaticID},
		Admins:      []primitive.ObjectID{userStaticID},
		Description: description,
		Owner:       &userStaticID,
	})

	return s.chatRepo.Create(ctx, *newChat)
}
