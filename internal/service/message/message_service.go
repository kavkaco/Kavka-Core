package message

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/utils/vali"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService interface {
	FetchMessages(ctx context.Context, chatID model.ChatID) ([]model.Message, *vali.Varror)
	UpdateTextMessage(ctx context.Context, chatID model.ChatID, newMessageContent string) *vali.Varror
	InsertTextMessage(ctx context.Context, chatID model.ChatID, userID model.UserID, messageContent string) (*model.Message, *vali.Varror)
	DeleteMessage(ctx context.Context, chatID model.ChatID, userID model.UserID, messageID model.MessageID) *vali.Varror
}

type MessageManager struct {
	messageRepo repository.MessageRepository
	chatRepo    repository.ChatRepository
	validator   *vali.Vali
}

func NewMessageService(messageRepo repository.MessageRepository, chatRepo repository.ChatRepository) MessageService {
	return &MessageManager{messageRepo, chatRepo, vali.Validator()}
}

func (s *MessageManager) FetchMessages(ctx context.Context, chatID model.ChatID) ([]model.Message, *vali.Varror) {
	messages, err := s.messageRepo.FetchMessages(ctx, chatID)
	if err != nil {
		return nil, &vali.Varror{Error: err}
	}

	return messages, nil
}

func (s *MessageManager) InsertTextMessage(ctx context.Context, chatID model.ChatID, userID model.UserID, messageContent string) (*model.Message, *vali.Varror) {
	validationErrors := s.validator.Validate(InsertTextMessageValidation{chatID, userID, messageContent})
	if len(validationErrors) > 0 {
		return nil, &vali.Varror{ValidationErrors: validationErrors}
	}

	chat, err := s.chatRepo.FindByID(ctx, chatID)
	if err != nil {
		return nil, &vali.Varror{Error: ErrChatNotFound}
	}

	if HasAccessToSendMessage(chat.ChatType, chat.ChatDetail, userID) {
		messageModel := model.NewMessage(model.TypeTextMessage, model.TextMessage{
			Text: messageContent,
		}, userID)
		message, err := s.messageRepo.Insert(ctx, chatID, messageModel)
		if err != nil {
			return nil, &vali.Varror{Error: ErrInsertMessage}
		}

		return message, nil
	}

	return nil, &vali.Varror{Error: ErrAccessDenied}
}

func (s *MessageManager) DeleteMessage(ctx context.Context, chatID model.ChatID, userID model.UserID, messageID model.MessageID) *vali.Varror {
	validationErrors := s.validator.Validate(DeleteMessageValidation{chatID, userID, messageID})
	if len(validationErrors) > 0 {
		return &vali.Varror{ValidationErrors: validationErrors}
	}

	chat, err := s.chatRepo.FindByID(ctx, chatID)
	if err != nil {
		return &vali.Varror{Error: ErrChatNotFound}
	}

	message, err := s.messageRepo.FindMessage(ctx, chatID, messageID)
	if err != nil {
		return &vali.Varror{Error: ErrNotFound}
	}

	if HasAccessToDeleteMessage(chat.ChatType, chat.ChatDetail, userID, *message) {
		err = s.messageRepo.Delete(ctx, chatID, messageID)
		if err != nil {
			return &vali.Varror{Error: ErrDeleteMessage}
		}

		return nil
	}

	return &vali.Varror{Error: ErrAccessDenied}
}

// TODO - Implement UpdateTextMessage Method For MessageService
func (s *MessageManager) UpdateTextMessage(ctx context.Context, chatID primitive.ObjectID, newMessageContent string) *vali.Varror {
	panic("unimplemented")
}
