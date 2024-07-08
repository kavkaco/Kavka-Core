package chat

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/utils/vali"
)

type ChatService interface {
	GetChat(ctx context.Context, chatID model.ChatID) (*model.Chat, *vali.Varror)
	GetUserChats(ctx context.Context, userID model.UserID) ([]model.Chat, *vali.Varror)
	CreateDirect(ctx context.Context, userID model.UserID, recipientUserID model.UserID) (*model.Chat, *vali.Varror)
	CreateGroup(ctx context.Context, userID model.UserID, title string, username string, description string) (*model.Chat, *vali.Varror)
	CreateChannel(ctx context.Context, userID model.UserID, title string, username string, description string) (*model.Chat, *vali.Varror)
}

type ChatManager struct {
	chatRepo  repository.ChatRepository
	userRepo  repository.UserRepository
	validator *vali.Vali
}

func NewChatService(chatRepo repository.ChatRepository, userRepo repository.UserRepository) ChatService {
	return &ChatManager{chatRepo, userRepo, vali.Validator()}
}

// find single chat with chat id
func (s *ChatManager) GetChat(ctx context.Context, chatID model.ChatID) (*model.Chat, *vali.Varror) {
	validationErrors := s.validator.Validate(GetChatValidation{chatID})
	if len(validationErrors) > 0 {
		return nil, &vali.Varror{ValidationErrors: validationErrors}
	}

	chat, err := s.chatRepo.FindByID(ctx, chatID)
	if err != nil {
		return nil, &vali.Varror{Error: ErrNotFound}
	}

	return chat, nil
}

// get the chats that belongs to user
func (s *ChatManager) GetUserChats(ctx context.Context, userID model.UserID) ([]model.Chat, *vali.Varror) {
	validationErrors := s.validator.Validate(GetUserChatsValidation{userID})
	if len(validationErrors) > 0 {
		return nil, &vali.Varror{ValidationErrors: validationErrors}

	}

	user, err := s.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, &vali.Varror{Error: ErrNotFound}
	}

	userChatsListIDs := user.ChatsListIDs

	userChats, err := s.chatRepo.FindManyByChatID(ctx, userChatsListIDs)
	if err != nil {
		return nil, &vali.Varror{Error: ErrGetUserChats}
	}

	return userChats, nil
}

func (s *ChatManager) CreateDirect(ctx context.Context, userID model.UserID, recipientUserID model.UserID) (*model.Chat, *vali.Varror) {
	validationErrors := s.validator.Validate(CreateDirectValidation{userID, recipientUserID})
	if len(validationErrors) > 0 {
		return nil, &vali.Varror{ValidationErrors: validationErrors}
	}

	sides := [2]model.UserID{userID, recipientUserID}

	// Check to do not be duplicated!
	dup, _ := s.chatRepo.FindBySides(ctx, sides)
	if dup != nil {
		return nil, &vali.Varror{Error: ErrChatAlreadyExists}
	}

	chatModel := model.NewChat(model.TypeDirect, &model.DirectChatDetail{
		Sides: sides,
	})

	saved, err := s.chatRepo.Create(ctx, *chatModel)
	if err != nil {
		return nil, &vali.Varror{Error: ErrCreateChat}
	}

	return saved, nil
}

func (s *ChatManager) CreateGroup(ctx context.Context, userID model.UserID, title string, username string, description string) (*model.Chat, *vali.Varror) {
	validationErrors := s.validator.Validate(CreateGroupValidation{userID, title, username, description})
	if len(validationErrors) > 0 {
		return nil, &vali.Varror{ValidationErrors: validationErrors}
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
		return nil, &vali.Varror{Error: ErrCreateChat}
	}

	return saved, nil
}

func (s *ChatManager) CreateChannel(ctx context.Context, userID model.UserID, title string, username string, description string) (*model.Chat, *vali.Varror) {
	validationErrors := s.validator.Validate(CreateChannelValidation{userID, title, username, description})
	if len(validationErrors) > 0 {
		return nil, &vali.Varror{ValidationErrors: validationErrors}
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
		return nil, &vali.Varror{Error: ErrCreateChat}
	}

	return saved, nil
}
