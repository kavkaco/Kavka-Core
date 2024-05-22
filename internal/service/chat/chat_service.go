package chat

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"go.uber.org/zap"
)

type ChatService interface {
	GetChat(ctx context.Context, chatID model.ChatID) (*model.Chat, error)
	GetUserChats(ctx context.Context, userID model.UserID) ([]model.Chat, error)
	CreateDirect(ctx context.Context, userID model.UserID, recipientUserID model.UserID) (*model.Chat, error)
	CreateGroup(ctx context.Context, userID model.UserID, title string, username string, description string) (*model.Chat, error)
	CreateChannel(ctx context.Context, userID model.UserID, title string, username string, description string) (*model.Chat, error)
}

type ChatManager struct {
	chatRepo repository.ChatRepository
	userRepo repository.UserRepository
}

func NewChatService(logger *zap.Logger, chatRepo repository.ChatRepository, userRepo repository.UserRepository) ChatService {
	return &ChatManager{chatRepo, userRepo}
}

// find single chat with chat id
func (s *ChatManager) GetChat(ctx context.Context, chatID model.ChatID) (*model.Chat, error) {
	chat, err := s.chatRepo.FindByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatManager) GetUserChats(ctx context.Context, userID model.UserID) ([]model.Chat, error) {
	user, err := s.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	userChatsListIDs := user.ChatsListIDs

	userChats, err := s.chatRepo.FindMany(ctx, userChatsListIDs)
	if err != nil {
		return nil, ErrGetUserChats
	}

	return userChats, nil
}

func (s *ChatManager) CreateDirect(ctx context.Context, userID model.UserID, recipientUserID model.UserID) (*model.Chat, error) {
	sides := [2]model.UserID{userID, recipientUserID}

	// Check to do not be duplicated!
	dup, _ := s.chatRepo.FindBySides(ctx, sides)
	if dup != nil {
		return nil, repository.ErrChatAlreadyExists
	}

	chatModel := model.NewChat(model.TypeDirect, &model.DirectChatDetail{
		Sides: sides,
	})

	saved, err := s.chatRepo.Create(ctx, *chatModel)
	if err != nil {
		return nil, ErrCreateChat
	}

	return saved, nil
}

func (s *ChatManager) CreateGroup(ctx context.Context, userID model.UserID, title string, username string, description string) (*model.Chat, error) {
	chatModel := model.NewChat(model.TypeGroup, &model.GroupChatDetail{
		Title:       title,
		Username:    username,
		Members:     []model.UserID{userID},
		Admins:      []model.UserID{userID},
		Description: description,
		Owner:       &userID,
	})

	saved, err := s.chatRepo.Create(ctx, *chatModel)
	if err != nil {
		return nil, ErrCreateChat
	}

	return saved, nil
}

func (s *ChatManager) CreateChannel(ctx context.Context, userID model.UserID, title string, username string, description string) (*model.Chat, error) {
	chatModel := model.NewChat(model.TypeGroup, &model.GroupChatDetail{
		Title:       title,
		Username:    username,
		Members:     []model.UserID{userID},
		Admins:      []model.UserID{userID},
		Description: description,
		Owner:       &userID,
	})

	saved, err := s.chatRepo.Create(ctx, *chatModel)
	if err != nil {
		return nil, ErrCreateChat
	}

	return saved, nil
}
