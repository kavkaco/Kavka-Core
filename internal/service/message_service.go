package service

import (
	"Kavka/internal/domain/chat"
	"Kavka/internal/domain/message"
	chatRepository "Kavka/internal/repository/chat"
	messageRepository "Kavka/internal/repository/message"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slices"
)

type MessageService struct {
	msgRepo  *messageRepository.MessageRepository
	chatRepo *chatRepository.ChatRepository
}

func NewMessageService(msgRepo *messageRepository.MessageRepository, chatRepo *chatRepository.ChatRepository) *MessageService {
	return &MessageService{msgRepo, chatRepo}
}

// This function is used to check that user can send message or not
func (s *MessageService) hasAccessToSendMessage(chatID primitive.ObjectID, staticID primitive.ObjectID) (bool, error) {
	foundChat, chatErr := s.chatRepo.FindByID(chatID)
	if chatErr != nil {
		return false, chatErr
	}

	if foundChat == nil {
		return false, chatRepository.ErrChatNotFound
	}

	log.Println(foundChat)

	if foundChat.ChatType == chat.ChatTypeDirect {
		hasSide := foundChat.ChatDetail.(*chat.DirectChatDetail).HasSide(staticID)
		return hasSide, nil
	} else if foundChat.ChatType == chat.ChatTypeChannel {
		var chatDetail chat.ChannelChatDetail

		chatDetailBSON, _ := chat.GetChatDetailBSON(foundChat.ChatDetail)
		bson.Unmarshal(chatDetailBSON, &chatDetail)

		admins := chatDetail.Admins
		isAdmin := slices.Contains(admins, &staticID)

		return isAdmin, nil
	} else if foundChat.ChatType == chat.ChatTypeGroup {
		var chatDetail chat.GroupChatDetail

		chatDetailBSON, _ := chat.GetChatDetailBSON(foundChat.ChatDetail)
		bson.Unmarshal(chatDetailBSON, &chatDetail)

		members := chatDetail.Members
		isMember := slices.Contains(members, &staticID)

		return isMember, nil
	}

	return true, nil
}

func (s *MessageService) InsertTextMessage(chatID primitive.ObjectID, staticID primitive.ObjectID, messageContent string) (*message.Message, error) {
	hasAccess, hasAccessErr := s.hasAccessToSendMessage(chatID, staticID)
	if hasAccessErr != nil {
		return nil, hasAccessErr
	}

	if hasAccess {
		msg := message.NewMessage(staticID, message.TypeTextMessage, &message.TextMessage{
			Message: messageContent,
		})

		return s.msgRepo.Insert(chatID, msg)
	}

	return nil, messageRepository.ErrNoAccess
}

func (s *MessageService) DeleteMessage(chatID primitive.ObjectID, messageID primitive.ObjectID) error {
	return s.msgRepo.Delete(chatID, messageID)
}
