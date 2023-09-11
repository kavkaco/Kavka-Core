package service

import (
	"Kavka/internal/domain/chat"
	chatRepository "Kavka/internal/repository/chat"
	userRepository "Kavka/internal/repository/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatService struct {
	chatRepo *chatRepository.ChatRepository
	userRepo *userRepository.UserRepository
}

func NewChatService(chatRepo *chatRepository.ChatRepository, userRepo *userRepository.UserRepository) *ChatService {
	return &ChatService{chatRepo, userRepo}
}

func (s *ChatService) GetChat(staticID primitive.ObjectID) (*chat.Chat, error) {
	foundChat, foundChatErr := s.chatRepo.FindChatOrSidesByStaticID(&staticID)
	if foundChatErr != nil {
		return nil, foundChatErr
	}

	return foundChat, nil
}

func (s *ChatService) CreateDirect(userStaticID primitive.ObjectID,
	targetStaticID primitive.ObjectID,
) (*chat.Chat, error) {
	sides := [2]*primitive.ObjectID{
		&userStaticID,
		&targetStaticID,
	}

	dup, _ := s.chatRepo.FindBySides(sides)
	if dup != nil {
		return nil, chatRepository.ErrChatAlreadyExists
	}

	return s.chatRepo.Create(chat.ChatTypeDirect, &chat.DirectChatDetail{
		Sides: sides,
	})
}

func (s *ChatService) CreateGroup(userStaticID primitive.ObjectID,
	title string, username string, description string,
) (*chat.Chat, error) {
	return s.chatRepo.Create(chat.ChatTypeGroup, &chat.GroupChatDetail{
		Title:       title,
		Username:    username,
		Members:     []*primitive.ObjectID{&userStaticID},
		Admins:      []*primitive.ObjectID{&userStaticID},
		Description: description,
	})
}

func (s *ChatService) CreateChannel(userStaticID primitive.ObjectID,
	title string, username string, description string,
) (*chat.Chat, error) {
	return s.chatRepo.Create(chat.ChatTypeGroup, &chat.GroupChatDetail{
		Title:       title,
		Username:    username,
		Members:     []*primitive.ObjectID{&userStaticID},
		Admins:      []*primitive.ObjectID{&userStaticID},
		Description: description,
	})
}
