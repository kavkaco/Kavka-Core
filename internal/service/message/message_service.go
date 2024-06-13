package message

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService interface {
	UpdateTextMessage(ctx context.Context, chatID model.ChatID, newMessageContent string) error
	InsertTextMessage(ctx context.Context, chatID model.ChatID, userID model.UserID, messageContent string) (*model.Message, error)
	DeleteMessage(ctx context.Context, chatID model.ChatID, userID model.UserID, messageID model.MessageID) error
}

type MessageManager struct {
	messageRepo repository.MessageRepository
	chatRepo    repository.ChatRepository
	validator   *validator.Validate
}

func NewMessageService(messageRepo repository.MessageRepository, chatRepo repository.ChatRepository) MessageService {
	validator := validator.New()
	return &MessageManager{messageRepo, chatRepo, validator}
}

func (s *MessageManager) InsertTextMessage(ctx context.Context, chatID model.ChatID, userID model.UserID, messageContent string) (*model.Message, error) {
	err := s.validator.Struct(InsertTextMessageValidation{chatID, userID, messageContent})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	chat, err := s.chatRepo.FindByID(ctx, chatID)
	if err != nil {
		return nil, ErrChatNotFound
	}

	if HasAccessToSendMessage(chat.ChatType, chat.ChatDetail, userID) {
		messageModel := model.NewMessage(model.TypeTextMessage, model.TextMessage{
			Data: messageContent,
		}, userID)
		message, err := s.messageRepo.Insert(ctx, chatID, messageModel)
		if err != nil {
			return nil, ErrInsertMessage
		}

		return message, nil
	}

	return nil, ErrAccessDenied
}

func (s *MessageManager) DeleteMessage(ctx context.Context, chatID model.ChatID, userID model.UserID, messageID model.MessageID) error {
	err := s.validator.Struct(DeleteMessageValidation{chatID, userID, messageID})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	chat, err := s.chatRepo.FindByID(ctx, chatID)
	if err != nil {
		return ErrChatNotFound
	}

	message, err := s.messageRepo.FindMessage(ctx, chatID, messageID)
	if err != nil {
		return ErrNotFound
	}

	if HasAccessToDeleteMessage(chat.ChatType, chat.ChatDetail, userID, *message) {
		err = s.messageRepo.Delete(ctx, chatID, messageID)
		if err != nil {
			return ErrDeleteMessage
		}

		return nil
	}

	return ErrAccessDenied
}

// UpdateTextMessage implements MessageService.
func (s *MessageManager) UpdateTextMessage(ctx context.Context, chatID primitive.ObjectID, newMessageContent string) error {
	panic("unimplemented")
}
