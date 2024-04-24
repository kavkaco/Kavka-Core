package service

import (
	"slices"

	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"github.com/kavkaco/Kavka-Core/internal/domain/message"
	messageRepository "github.com/kavkaco/Kavka-Core/internal/repository/message"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type messageService struct {
	logger   *zap.Logger
	msgRepo  message.Repository
	chatRepo chat.Repository
}

func NewMessageService(logger *zap.Logger, msgRepo message.Repository, chatRepo chat.Repository) message.Service {
	return &messageService{logger, msgRepo, chatRepo}
}

// This function is used to check that user can send message or not.
func (s *messageService) hasAccessToSendMessage(chatID primitive.ObjectID, staticID primitive.ObjectID) (bool, error) {
	foundChat, chatErr := s.chatRepo.FindByID(chatID)
	if chatErr != nil {
		return false, chatErr
	}

	if foundChat.ChatType == chat.TypeDirect {
		hasSide := foundChat.ChatDetail.(*chat.DirectChatDetail).HasSide(staticID)
		return hasSide, nil
	} else if foundChat.ChatType == chat.TypeChannel {
		admins := foundChat.ChatDetail.(*chat.ChannelChatDetail).Admins
		isAdmin := slices.Contains(admins, staticID)
		return isAdmin, nil
	} else if foundChat.ChatType == chat.TypeGroup {
		// log.Println(reflect.TypeOf(foundChat.ChatDetail).Name())
		// members := foundChat.ChatDetail.(*chat.GroupChatDetail).Members
		// isMember := slices.Contains(members, &staticID)
		// return isMember, nil
		// FIXME
		// log.Println(foundChat.(chat.Chat).ChatDetail.(chat.GroupChatDetail))
		return true, nil
	}

	return false, nil
}

// This function is used to check that user can delete message or not.
// func (s *messageService) hasAccessToDeleteMessage(chatID primitive.ObjectID, staticID primitive.ObjectID) (bool, error) {
// 	foundChat, chatErr := s.chatRepo.FindByID(chatID)
// 	if chatErr != nil {
// 		return false, chatErr
// 	}

// 	if foundChat.ChatType == chat.ChatTypeDirect {
// 		hasSide := foundChat.ChatDetail.(*chat.DirectChatDetail).HasSide(staticID)
// 		return hasSide, nil
// 	}

// 	if foundChat.ChatType == chat.ChatTypeChannel {
// 		admins := foundChat.ChatDetail.(chat.ChannelChatDetail).Admins
// 		isAdmin := slices.Contains(admins, &staticID)
// 		return isAdmin, nil
// 	}

// 	if foundChat.ChatType == chat.ChatTypeGroup {
// 		members := foundChat.ChatDetail.(chat.GroupChatDetail).Members
// 		isMember := slices.Contains(members, &staticID)
// 		return isMember, nil
// 	}

// 	return false, nil
// }

func (s *messageService) InsertTextMessage(chatID primitive.ObjectID, staticID primitive.ObjectID, messageContent string) (*message.Message, error) {
	hasAccess, hasAccessErr := s.hasAccessToSendMessage(chatID, staticID)
	if hasAccessErr != nil {
		return nil, hasAccessErr
	}

	if hasAccess {
		msg := message.NewMessage(staticID, message.TypeTextMessage, &message.TextMessage{
			Data: messageContent,
		})

		return s.msgRepo.Insert(chatID, msg)
	}

	return nil, messageRepository.ErrNoAccess
}

func (s *messageService) DeleteMessage(chatID primitive.ObjectID, messageID primitive.ObjectID) error {
	return s.msgRepo.Delete(chatID, messageID)
}
