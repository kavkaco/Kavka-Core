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
	useV2          bool
}

type MessageServiceOption func(*MessageService)

func WithV2Schema(useV2 bool) MessageServiceOption {
	return func(s *MessageService) {
		s.useV2 = useV2
	}
}

func NewMessageService(logger *log.SubLogger, messageRepo repository.MessageRepository, chatRepo repository.ChatRepository, userRepo repository.UserRepository, eventPublisher stream.StreamPublisher, opts ...MessageServiceOption) *MessageService {
	s := &MessageService{logger, messageRepo, chatRepo, userRepo, vali.Validator(), eventPublisher, false}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *MessageService) FetchMessages(ctx context.Context, chatID model.ChatID) ([]*model.MessageGetter, *vali.ValiErr) {
	return s.FetchMessagesPaginated(ctx, chatID, 0, DefaultMessagesPerPage)
}

func (s *MessageService) FetchMessagesPaginated(ctx context.Context, chatID model.ChatID, skip, limit int64) ([]*model.MessageGetter, *vali.ValiErr) {
	if limit <= 0 || limit > DefaultMessagesPerPage {
		limit = DefaultMessagesPerPage
	}

	if s.useV2 {
		docs, err := s.messageRepo.FetchMessagesV2(ctx, chatID, skip, limit)
		if err != nil {
			return nil, &vali.ValiErr{Error: err}
		}

		result := make([]*model.MessageGetter, 0, len(docs))
		for _, doc := range docs {
			result = append(result, doc.ToMessageGetter())
		}
		return result, nil
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

	u, err := s.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, &vali.ValiErr{Error: err}
	}

	var messageGetter *model.MessageGetter

	if s.useV2 {
		doc, err := s.messageRepo.InsertV2(ctx, &model.MessageDocumentV2{
			ChatID:         chatID,
			SenderID:       userID,
			SenderName:     u.Name,
			SenderLastName: u.LastName,
			SenderUsername: u.Username,
			Type:           model.TypeTextMessage,
			Content:        model.TextMessage{Text: messageContent},
		})
		if err != nil {
			return nil, &vali.ValiErr{Error: ErrInsertMessage}
		}

		messageGetter = doc.ToMessageGetter()
	} else {
		m, err := s.messageRepo.Insert(ctx, chatID, model.NewMessage(model.TypeTextMessage, model.TextMessage{
			Text: messageContent,
		}, userID))
		if err != nil {
			return nil, &vali.ValiErr{Error: ErrInsertMessage}
		}

		messageGetter = &model.MessageGetter{
			Sender: &model.MessageSenderDTO{
				UserID:   u.UserID,
				Name:     u.Name,
				LastName: u.LastName,
				Username: u.Username,
			},
			Message: m,
		}
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

	if s.useV2 {
		doc, err := s.messageRepo.FetchMessageV2(ctx, chatID, messageID)
		if err != nil {
			return &vali.ValiErr{Error: ErrNotFound}
		}

		msg := &model.Message{
			MessageID: doc.ID,
			SenderID:  doc.SenderID,
		}

		if HasAccessToDeleteMessage(chat.ChatType, chat.ChatDetail, userID, *msg) {
			err = s.messageRepo.DeleteV2(ctx, chatID, messageID)
			if err != nil {
				return &vali.ValiErr{Error: ErrDeleteMessage}
			}

			return nil
		}

		return &vali.ValiErr{Error: ErrAccessDenied}
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

	if s.useV2 {
		doc, err := s.messageRepo.FetchMessageV2(ctx, chatID, messageID)
		if err != nil {
			return &vali.ValiErr{Error: ErrNotFound}
		}

		if doc.SenderID != userID {
			return &vali.ValiErr{Error: ErrAccessDenied}
		}

		err = s.messageRepo.UpdateMessageContentV2(ctx, chatID, messageID, newMessageContent)
		if err != nil {
			return &vali.ValiErr{Error: err}
		}

		return nil
	}

	chat, err := s.chatRepo.GetChat(ctx, chatID)
	if err != nil {
		return &vali.ValiErr{Error: ErrChatNotFound}
	}

	_ = chat

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
