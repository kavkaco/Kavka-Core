package service

import (
	"Kavka/internal/domain/chat"
	"Kavka/internal/domain/message"
	chatRepository "Kavka/internal/repository/chat"
	messageRepository "Kavka/internal/repository/message"
	"Kavka/utils/slices"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	if foundChat.ChatType == chat.ChatTypeDirect {
		hasSide := foundChat.ChatDetail.(*chat.DirectChatDetail).HasSide(staticID)
		return hasSide, nil
	} else if foundChat.ChatType == chat.ChatTypeChannel {
		var chatDetail chat.ChannelChatDetail

		chatDetailBSON, _ := chat.GetChatDetailBSON(foundChat.ChatDetail)
		bson.Unmarshal(chatDetailBSON, &chatDetail)

		admins := chatDetail.Admins
		isAdmin := slices.ContainsObjectID(admins, staticID)

		return isAdmin, nil
	} else if foundChat.ChatType == chat.ChatTypeGroup {
		var chatDetail chat.GroupChatDetail

		chatDetailBSON, _ := chat.GetChatDetailBSON(foundChat.ChatDetail)
		bson.Unmarshal(chatDetailBSON, &chatDetail)

		members := chatDetail.Members
		isMember := slices.ContainsObjectID(members, staticID)

		return isMember, nil
	}

	return true, nil
}

// This function is used to check that user can delete message or not
func (s *MessageService) hasAccessToDeleteMessage(chatID primitive.ObjectID, staticID primitive.ObjectID, messageID primitive.ObjectID) (bool, error) {
	foundChat, chatErr := s.chatRepo.FindByID(chatID)
	if chatErr != nil {
		return false, chatErr
	}

	if foundChat == nil {
		return false, chatRepository.ErrChatNotFound
	}

	if foundChat.ChatType == chat.ChatTypeDirect {
		// This means every side of the chat can delete the message for both sides
		hasSide := foundChat.ChatDetail.(*chat.DirectChatDetail).HasSide(staticID)
		return hasSide, nil
	} else if foundChat.ChatType == chat.ChatTypeChannel {
		// Only an admin can remove a message from the chat
		var chatDetail chat.ChannelChatDetail

		chatDetailBSON, _ := chat.GetChatDetailBSON(foundChat.ChatDetail)
		bson.Unmarshal(chatDetailBSON, &chatDetail)

		admins := chatDetail.Admins
		isAdmin := slices.ContainsObjectID(admins, staticID)

		return isAdmin, nil
	} else if foundChat.ChatType == chat.ChatTypeGroup {
		// A user only can delete her own message
		var chatDetail chat.GroupChatDetail

		chatDetailBSON, _ := chat.GetChatDetailBSON(foundChat.ChatDetail)
		bson.Unmarshal(chatDetailBSON, &chatDetail)

		members := chatDetail.Members
		isMember := slices.ContainsObjectID(members, staticID)

		// No body could be able to make a change in a non-joined chat
		if !isMember {
			return false, messageRepository.ErrIsNotAMember
		}

		msg := foundChat.GetMessage(messageID)
		if msg == nil {
			// No message to delete
			return false, messageRepository.ErrMessageNotFound
		}

		// Check to be user's own message
		if msg.SenderID.Hex() != staticID.Hex() {
			return true, messageRepository.ErrIsNotUsersOwnMessage
		}
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

func (s *MessageService) DeleteMessage(chatID primitive.ObjectID, staticID primitive.ObjectID, messageID primitive.ObjectID) error {
	hasAccess, hasAccessErr := s.hasAccessToDeleteMessage(chatID, staticID, messageID)
	if hasAccessErr != nil {
		return hasAccessErr
	}

	if hasAccess {
		return s.msgRepo.Delete(chatID, messageID)
	}

	return messageRepository.ErrNoAccess
}
