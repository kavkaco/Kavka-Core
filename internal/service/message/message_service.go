package message

import (
	"context"

	"github.com/kavkaco/Kavka-Core/infra/stream"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/model/proto_model_transformer"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/log"
	"github.com/kavkaco/Kavka-Core/utils/vali"
	eventsv1 "github.com/kavkaco/Kavka-ProtoBuf/gen/go/protobuf/events/v1"
	"google.golang.org/protobuf/proto"
)

const DefaultMessagesPerPage = 200

type MessageService struct {
	logger         *log.SubLogger
	messageRepo    repository.MessageRepository
	chatRepo       repository.ChatRepository
	userRepo       repository.UserRepository
	validator      *vali.Vali
	eventPublisher stream.StreamPublisher
}

func NewMessageService(logger *log.SubLogger, messageRepo repository.MessageRepository, chatRepo repository.ChatRepository, userRepo repository.UserRepository, eventPublisher stream.StreamPublisher) *MessageService {
	return &MessageService{logger, messageRepo, chatRepo, userRepo, vali.Validator(), eventPublisher}
}

func (s *MessageService) FetchMessages(ctx context.Context, chatID model.ChatID) ([]*model.MessageGetter, *vali.ValiErr) {
	return s.FetchMessagesPaginated(ctx, chatID, 0, DefaultMessagesPerPage)
}

func (s *MessageService) FetchMessagesPaginated(ctx context.Context, chatID model.ChatID, skip, limit int64) ([]*model.MessageGetter, *vali.ValiErr) {
	if limit <= 0 || limit > DefaultMessagesPerPage {
		limit = DefaultMessagesPerPage
	}

	messages, err := s.messageRepo.FetchMessagesPaginated(ctx, chatID, skip, limit)
	if err != nil {
		return nil, &vali.ValiErr{Error: err}
	}

	return messages, nil
}

func (s *MessageService) SendTextMessage(ctx context.Context, chatID model.ChatID, userID model.UserID, messageContent string) (*model.MessageGetter, *vali.ValiErr) {
	errs := s.validator.Validate(insertTextMessageValidation{chatID, userID, messageContent})
	if len(errs) > 0 {
		return nil, &vali.ValiErr{ValidationErrors: errs}
	}

	c, err := s.chatRepo.GetChat(ctx, chatID)
	if err != nil {
		return nil, &vali.ValiErr{Error: ErrChatNotFound}
	}

	if !HasAccessToSendMessage(c.ChatType, c.ChatDetail, userID) {
		return nil, &vali.ValiErr{Error: ErrAccessDenied}
	}

	m, err := s.messageRepo.Insert(ctx, chatID, model.NewMessage(model.TypeTextMessage, model.TextMessage{
		Text: messageContent,
	}, userID))
	if err != nil {
		return nil, &vali.ValiErr{Error: ErrInsertMessage}
	}

	u, err := s.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, &vali.ValiErr{Error: err}
	}

	messageGetter := &model.MessageGetter{
		Sender: &model.MessageSenderDTO{
			UserID:   u.UserID,
			Name:     u.Name,
			LastName: u.LastName,
			Username: u.Username,
		},
		Message: m,
	}

	go func() {
		eventReceivers, receiversErr := ReceiversIDs(c)
		if receiversErr != nil {
			s.logger.Error(receiversErr.Error())
			return
		}

		payloadProtoBuf, marshalErr := proto.Marshal(&eventsv1.SubscribeEventsStreamResponse{
			Name: "add-message",
			Type: eventsv1.SubscribeEventsStreamResponse_TYPE_ADD_MESSAGE,
			Payload: &eventsv1.SubscribeEventsStreamResponse_AddMessage{
				AddMessage: &eventsv1.AddMessage{
					ChatId:  chatID.Hex(),
					Message: proto_model_transformer.MessageToProto(messageGetter),
				},
			},
		},
		)
		if marshalErr != nil {
			s.logger.Error("proto marshal error: " + marshalErr.Error())
			return
		}

		publishErr := s.eventPublisher.Publish(&eventsv1.StreamEvent{
			SenderUserId:    userID,
			ReceiversUserId: eventReceivers,
			Payload:         payloadProtoBuf,
		})
		if publishErr != nil {
			s.logger.Error("unable to publish add-chat event in eventPublisher: " + publishErr.Error())
		}
	}()

	return messageGetter, nil
}

func (s *MessageService) DeleteMessage(ctx context.Context, chatID model.ChatID, userID model.UserID, messageID model.MessageID) *vali.ValiErr {
	errs := s.validator.Validate(deleteMessageValidation{chatID, userID, messageID})
	if len(errs) > 0 {
		return &vali.ValiErr{ValidationErrors: errs}
	}

	chat, err := s.chatRepo.GetChat(ctx, chatID)
	if err != nil {
		return &vali.ValiErr{Error: ErrChatNotFound}
	}

	message, err := s.messageRepo.FetchMessage(ctx, chatID, messageID)
	if err != nil {
		return &vali.ValiErr{Error: ErrNotFound}
	}

	if HasAccessToDeleteMessage(chat.ChatType, chat.ChatDetail, userID, *message) {
		err = s.messageRepo.Delete(ctx, chatID, messageID)
		if err != nil {
			return &vali.ValiErr{Error: ErrDeleteMessage}
		}

		return nil
	}

	return &vali.ValiErr{Error: ErrAccessDenied}
}

func (s *MessageService) UpdateTextMessage(ctx context.Context, chatID model.ChatID, userID model.UserID, messageID model.MessageID, newMessageContent string) *vali.ValiErr {
	errs := s.validator.Validate(updateTextMessageValidation{chatID, userID, messageID, newMessageContent})
	if len(errs) > 0 {
		return &vali.ValiErr{ValidationErrors: errs}
	}

	chat, err := s.chatRepo.GetChat(ctx, chatID)
	if err != nil {
		return &vali.ValiErr{Error: ErrChatNotFound}
	}

	message, err := s.messageRepo.FetchMessage(ctx, chatID, messageID)
	if err != nil {
		return &vali.ValiErr{Error: ErrNotFound}
	}

	if message.SenderID != userID {
		return &vali.ValiErr{Error: ErrAccessDenied}
	}

	err = s.messageRepo.UpdateMessageContent(ctx, chatID, messageID, newMessageContent)
	if err != nil {
		return &vali.ValiErr{Error: err}
	}

	return nil
}
