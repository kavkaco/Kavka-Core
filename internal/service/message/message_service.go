package service

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/model/chat"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type MessageService interface {
	InsertTextMessage(ctx context.Context, chatID primitive.ObjectID, staticID primitive.ObjectID, messageContent string) (*model.Message, error)
	DeleteMessage(ctx context.Context, chatID primitive.ObjectID, messageID primitive.ObjectID) error
}

type MessageManager struct {
	logger   *zap.Logger
	msgRepo  repository.MessageRepository
	chatRepo repository.ChatRepository
}

func NewMessageService(logger *zap.Logger, msgRepo repository.MessageRepository, chatRepo repository.ChatRepository) MessageService {
	return &MessageManager{logger, msgRepo, chatRepo}
}

// This function is used to find the chat and check that user has access to send message or not.
func (s *MessageManager) hasAccessToSendMessage(ctx context.Context, chatID primitive.ObjectID, staticID primitive.ObjectID) (bool, error) {
	foundChat, chatErr := s.chatRepo.FindByID(ctx, chatID)
	if chatErr != nil {
		return false, chatErr
	}

	if foundChat.ChatType == chat.TypeDirect {
		hasSide := foundChat.ChatDetail.(*chat.DirectChatDetail).HasSide(staticID)
		return hasSide, nil
	} else if foundChat.ChatType == chat.TypeChannel {
		chatDetail, err := utils.TypeConverter[chat.ChannelChatDetail](foundChat.ChatDetail)
		if err != nil {
			return false, err
		}

		return chatDetail.HasAccessToSendMessage(staticID), nil
	} else if foundChat.ChatType == chat.TypeGroup {
		chatDetail, err := utils.TypeConverter[chat.GroupChatDetail](foundChat.ChatDetail)
		if err != nil {
			return false, err
		}

		return chatDetail.HasAccessToSendMessage(staticID), nil
	}

	return false, nil
}

// This function is used to find the chat and check that user has access to send message or not.
func (s *MessageManager) hasAccessToDeleteMessage(ctx context.Context, chatID primitive.ObjectID, staticID primitive.ObjectID) (bool, error) {
	foundChat, chatErr := s.chatRepo.FindByID(ctx, chatID)
	if chatErr != nil {
		return false, chatErr
	}

	if foundChat.ChatType == chat.TypeDirect {
		hasSide := foundChat.ChatDetail.(*chat.DirectChatDetail).HasSide(staticID)
		return hasSide, nil
	} else if foundChat.ChatType == chat.TypeChannel {
		chatDetail, err := utils.TypeConverter[chat.ChannelChatDetail](foundChat.ChatDetail)
		if err != nil {
			return false, err
		}

		return chatDetail.HasAccessToDeleteMessage(staticID), nil
	} else if foundChat.ChatType == chat.TypeGroup {
		// Get the message from messages_collection
		chatDetail, err := utils.TypeConverter[chat.GroupChatDetail](foundChat.ChatDetail)
		if err != nil {
			return false, err
		}

		// FIXME
		println(chatDetail)
		// return chatDetail.HasAccessToDeleteMessage(staticID, msg), nil
	}

	return false, nil
}

func (s *MessageManager) InsertTextMessage(ctx context.Context, chatID primitive.ObjectID, staticID primitive.ObjectID, messageContent string) (*model.Message, error) {
	// hasAccess, hasAccessErr := s.hasAccessToSendMessage(ctx, chatID, staticID)
	// if hasAccessErr != nil {
	// 	return nil, hasAccessErr
	// }

	// if hasAccess {
	// 	msg := model.NewMessage(staticID, model.TypeTextMessage, &model.TextMessage{
	// 		Data: messageContent,
	// 	})

	// 	return s.msgRepo.Insert(ctx, chatID, msg)
	// }

	// return nil, messageRepository.ErrNoAccess
	return nil, nil
}

func (s *MessageManager) DeleteMessage(ctx context.Context, chatID primitive.ObjectID, messageID primitive.ObjectID) error {
	return s.msgRepo.Delete(ctx, chatID, messageID)
}
