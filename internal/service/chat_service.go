package service

import (
	"github.com/kavkaco/Kavka-Core/internal/model/chat"
	"github.com/kavkaco/Kavka-Core/internal/model/user"
	chatRepository "github.com/kavkaco/Kavka-Core/internal/repository/chat"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type chatService struct {
	logger   *zap.Logger
	chatRepo chat.Repository
	userRepo user.UserRepository
}

func NewChatService(logger *zap.Logger, chatRepo chat.Repository, userRepo user.UserRepository) chat.Service {
	return &chatService{logger, chatRepo, userRepo}
}

func (s *chatService) GetChat(staticID primitive.ObjectID) (*chat.ChatC, error) {
	foundChat, err := s.chatRepo.FindChatOrSidesByStaticID(staticID)
	if err != nil {
		return nil, err
	}

	return foundChat, nil
}

func (s *chatService) GetUserChats(userStaticID primitive.ObjectID) ([]chat.ChatC, error) {
	return []chat.ChatC{}, nil
	// return s.chatRepo.GetChats(userStaticID)
}

func (s *chatService) CreateDirect(userStaticID primitive.ObjectID, targetStaticID primitive.ObjectID) (*chat.Chat, error) {
	sides := [2]primitive.ObjectID{
		userStaticID,
		targetStaticID,
	}

	// Check to do not be duplicated!
	dup, _ := s.chatRepo.FindBySides(sides)
	if dup != nil {
		return nil, chatRepository.ErrChatAlreadyExists
	}

	newChat := chat.NewChat(chat.TypeDirect, &chat.DirectChatDetail{
		Sides: sides,
	})

	return s.chatRepo.Create(*newChat)
}

func (s *chatService) CreateGroup(userStaticID primitive.ObjectID, title string, username string, description string) (*chat.Chat, error) {
	newChat := chat.NewChat(chat.TypeGroup, &chat.GroupChatDetail{
		Title:       title,
		Username:    username,
		Members:     []primitive.ObjectID{userStaticID},
		Admins:      []primitive.ObjectID{userStaticID},
		Description: description,
		Owner:       &userStaticID,
	})

	return s.chatRepo.Create(*newChat)
}

func (s *chatService) CreateChannel(userStaticID primitive.ObjectID, title string, username string, description string) (*chat.Chat, error) {
	newChat := chat.NewChat(chat.TypeGroup, &chat.GroupChatDetail{
		Title:       title,
		Username:    username,
		Members:     []primitive.ObjectID{userStaticID},
		Admins:      []primitive.ObjectID{userStaticID},
		Description: description,
		Owner:       &userStaticID,
	})

	return s.chatRepo.Create(*newChat)
}
