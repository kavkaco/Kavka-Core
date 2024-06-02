package chat

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
)

type ChatService interface {
	GetChat(ctx context.Context, chatID model.ChatID) (*model.Chat, error)
	GetUserChats(ctx context.Context, userID model.UserID) ([]model.Chat, error)
	CreateDirect(ctx context.Context, userID model.UserID, recipientUserID model.UserID) (*model.Chat, error)
	CreateGroup(ctx context.Context, userID model.UserID, title string, username string, description string) (*model.Chat, error)
	CreateChannel(ctx context.Context, userID model.UserID, title string, username string, description string) (*model.Chat, error)
}

type ChatManager struct {
	chatRepo  repository.ChatRepository
	userRepo  repository.UserRepository
	validator *validator.Validate
}

func NewChatService(chatRepo repository.ChatRepository, userRepo repository.UserRepository) ChatService {
	validator := validator.New()
	return &ChatManager{chatRepo, userRepo, validator}
}

// find single chat with chat id
func (s *ChatManager) GetChat(ctx context.Context, chatID model.ChatID) (*model.Chat, error) {
	err := s.validator.Struct(GetChatValidation{chatID})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	chat, err := s.chatRepo.FindByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

// get the chats that belongs to user
func (s *ChatManager) GetUserChats(ctx context.Context, userID model.UserID) ([]model.Chat, error) {
	err := s.validator.Struct(GetUserChatsValidation{userID})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	user, err := s.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	userChatsListIDs := user.ChatsListIDs

	userChats, err := s.chatRepo.FindManyByChatID(ctx, userChatsListIDs)
	if err != nil {
		return nil, ErrGetUserChats
	}

	return userChats, nil
}

func (s *ChatManager) CreateDirect(ctx context.Context, userID model.UserID, recipientUserID model.UserID) (*model.Chat, error) {
	err := s.validator.Struct(CreateDirectValidation{userID, recipientUserID})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

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
	err := s.validator.Struct(CreateGroupValidation{userID, title, username, description})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	chatModel := model.NewChat(model.TypeGroup, &model.GroupChatDetail{
		Title:       title,
		Username:    username,
		Members:     []model.UserID{userID},
		Admins:      []model.UserID{userID},
		Description: description,
		Owner:       userID,
	})

	saved, err := s.chatRepo.Create(ctx, *chatModel)
	if err != nil {
		return nil, ErrCreateChat
	}

	return saved, nil
}

func (s *ChatManager) CreateChannel(ctx context.Context, userID model.UserID, title string, username string, description string) (*model.Chat, error) {
	err := s.validator.Struct(CreateChannelValidation{userID, title, username, description})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	chatModel := model.NewChat(model.TypeChannel, &model.ChannelChatDetail{
		Title:       title,
		Username:    username,
		Members:     []model.UserID{userID},
		Admins:      []model.UserID{userID},
		Description: description,
		Owner:       userID,
	})

	saved, err := s.chatRepo.Create(ctx, *chatModel)
	if err != nil {
		return nil, ErrCreateChat
	}

	return saved, nil
}
