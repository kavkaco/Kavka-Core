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

// The function `GetOrCreateChat` checks the chat type and username, and if the chat type is direct, it gets or creates a chat.
// but if chat type was not direct and was group or channel we just can get that from repo and deliver it to user.
// this means that we can't create a group or channel in this function and it should be handled on another function.
func (s *ChatService) GetOrCreateChat(chatType string, chatUsername string, userStaticID primitive.ObjectID) (*chat.Chat, error) {
	if chatType == chat.ChatTypeDirect {
		foundUser, foundUserErr := s.userRepo.FindByUsername(chatUsername)
		if foundUserErr != nil {
			return nil, foundUserErr
		}

		foundChat, foundChatErr := s.chatRepo.FindBySides([2]*primitive.ObjectID{
			&userStaticID,
			&foundUser.StaticID,
		})
		if foundChatErr == chatRepository.ErrChatNotFound {
			// Create a new direct chat
			return s.chatRepo.Create(chat.ChatTypeDirect, &chat.DirectChatDetail{
				Sides: [2]*primitive.ObjectID{&foundUser.StaticID, &userStaticID},
			})
		} else if foundChatErr == nil {
			// Chat already exists and there is no needed to create a new one
			return foundChat, nil
		} else {
			return nil, foundChatErr
		}
	}

	return nil, nil
}
