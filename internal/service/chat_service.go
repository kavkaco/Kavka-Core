package service

import (
	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"github.com/kavkaco/Kavka-Core/internal/domain/user"
	chatRepository "github.com/kavkaco/Kavka-Core/internal/repository/chat"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type chatService struct {
	chatRepo chat.Repository
	userRepo user.UserRepository
}

func NewChatService(chatRepo chat.Repository, userRepo user.UserRepository) chat.Service {
	return &chatService{chatRepo, userRepo}
}

func (s *chatService) GetChat(staticID primitive.ObjectID) (*chat.Chat, error) {
	foundChat, foundChatErr := s.chatRepo.FindChatOrSidesByStaticID(staticID)
	if foundChatErr != nil {
		return nil, foundChatErr
	}

	return foundChat, nil
}

func (s *chatService) CreateDirect(userStaticID primitive.ObjectID, targetStaticID primitive.ObjectID) (*chat.Chat, error) {
	sides := [2]primitive.ObjectID{
		userStaticID,
		targetStaticID,
	}

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
	})

	return s.chatRepo.Create(*newChat)
}
