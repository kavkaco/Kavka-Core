package chat

import (
	"context"

	grpc_model "github.com/kavkaco/Kavka-Core/delivery/grpc/model"
	"github.com/kavkaco/Kavka-Core/infra/stream"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/log"
	eventsv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/events/v1"
	"github.com/kavkaco/Kavka-Core/utils/vali"
	"google.golang.org/protobuf/proto"
)

const SubjChats = "chats"

type ChatService interface {
	GetChat(ctx context.Context, chatID model.ChatID) (*model.Chat, *vali.Varror)
	GetUserChats(ctx context.Context, userID model.UserID) ([]model.Chat, *vali.Varror)
	CreateDirect(ctx context.Context, userID model.UserID, recipientUserID model.UserID) (*model.Chat, *vali.Varror)
	CreateGroup(ctx context.Context, userID model.UserID, title string, username string, description string) (*model.Chat, *vali.Varror)
	CreateChannel(ctx context.Context, userID model.UserID, title string, username string, description string) (*model.Chat, *vali.Varror)
}

type ChatManager struct {
	logger         *log.SubLogger
	chatRepo       repository.ChatRepository
	userRepo       repository.UserRepository
	validator      *vali.Vali
	eventPublisher stream.StreamPublisher
}

func NewChatService(logger *log.SubLogger, chatRepo repository.ChatRepository, userRepo repository.UserRepository, eventPublisher stream.StreamPublisher) ChatService {
	return &ChatManager{logger, chatRepo, userRepo, vali.Validator(), eventPublisher}
}

// find single chat with chat id
func (s *ChatManager) GetChat(ctx context.Context, chatID model.ChatID) (*model.Chat, *vali.Varror) {
	validationErrors := s.validator.Validate(GetChatValidation{chatID})
	if len(validationErrors) > 0 {
		return nil, &vali.Varror{ValidationErrors: validationErrors}
	}

	chat, err := s.chatRepo.FindByID(ctx, chatID)
	if err != nil {
		return nil, &vali.Varror{Error: ErrNotFound}
	}

	return chat, nil
}

// get the chats that belongs to user
func (s *ChatManager) GetUserChats(ctx context.Context, userID model.UserID) ([]model.Chat, *vali.Varror) {
	validationErrors := s.validator.Validate(GetUserChatsValidation{userID})
	if len(validationErrors) > 0 {
		return nil, &vali.Varror{ValidationErrors: validationErrors}
	}

	user, err := s.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, &vali.Varror{Error: ErrNotFound}
	}

	userChatsListIDs := user.ChatsListIDs

	userChats, err := s.chatRepo.FindManyByChatID(ctx, userChatsListIDs)
	if err != nil {
		return nil, &vali.Varror{Error: ErrGetUserChats}
	}

	return userChats, nil
}

func (s *ChatManager) CreateDirect(ctx context.Context, userID model.UserID, recipientUserID model.UserID) (*model.Chat, *vali.Varror) {
	validationErrors := s.validator.Validate(CreateDirectValidation{userID, recipientUserID})
	if len(validationErrors) > 0 {
		return nil, &vali.Varror{ValidationErrors: validationErrors}
	}

	sides := [2]model.UserID{userID, recipientUserID}

	// Check to do not be duplicated!
	dup, _ := s.chatRepo.FindBySides(ctx, sides)
	if dup != nil {
		return nil, &vali.Varror{Error: ErrChatAlreadyExists}
	}

	chatModel := model.NewChat(model.TypeDirect, &model.DirectChatDetail{
		Sides: sides,
	})

	saved, err := s.chatRepo.Create(ctx, *chatModel)
	if err != nil {
		return nil, &vali.Varror{Error: ErrCreateChat}
	}

	return saved, nil
}

func (s *ChatManager) CreateGroup(ctx context.Context, userID model.UserID, title string, username string, description string) (*model.Chat, *vali.Varror) {
	validationErrors := s.validator.Validate(CreateGroupValidation{userID, title, username, description})
	if len(validationErrors) > 0 {
		return nil, &vali.Varror{ValidationErrors: validationErrors}
	}

	chatModel := model.NewChat(model.TypeGroup, &model.GroupChatDetail{
		Title:       title,
		Username:    username,
		Members:     []model.UserID{userID},
		Admins:      []model.UserID{userID},
		Description: description,
		Owner:       userID,
	})

	saved, err := s.chatRepo.Create(ctx, *chatModel)
	if err != nil {
		return nil, &vali.Varror{Error: ErrCreateChat}
	}

	return saved, nil
}

func (s *ChatManager) CreateChannel(ctx context.Context, userID model.UserID, title string, username string, description string) (*model.Chat, *vali.Varror) {
	validationErrors := s.validator.Validate(CreateChannelValidation{userID, title, username, description})
	if len(validationErrors) > 0 {
		return nil, &vali.Varror{ValidationErrors: validationErrors}
	}

	chatModel := model.NewChat(model.TypeChannel, &model.ChannelChatDetail{
		Title:       title,
		Username:    username,
		Members:     []model.UserID{userID},
		Admins:      []model.UserID{userID},
		Description: description,
		Owner:       userID,
	})

	chatGrpcModel, err := grpc_model.TransformChatToGrpcModel(*chatModel)
	if err != nil {
		return nil, &vali.Varror{Error: grpc_model.ErrTransformToGrpcModel}
	}

	if s.eventPublisher != nil {
		go func() {
			payloadProtoBuf, marshalErr := proto.Marshal(&eventsv1.SubscribeEventsStreamResponse{
				Name: "add-chat",
				Type: eventsv1.SubscribeEventsStreamResponse_TYPE_ADD_CHAT,
				Payload: &eventsv1.SubscribeEventsStreamResponse_AddChat{
					AddChat: &eventsv1.AddChat{
						Chat: chatGrpcModel,
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
				ReceiversUserId: []string{userID},
				Payload:         payloadProtoBuf,
			})
			if publishErr != nil {
				s.logger.Error("unable to publish add-chat event in eventPublisher: " + publishErr.Error())
			}
		}()
	}

	saved, err := s.chatRepo.Create(ctx, *chatModel)
	if err != nil {
		return nil, &vali.Varror{Error: ErrCreateChat}
	}

	err = s.chatRepo.AddToUsersChatsList(ctx, userID, saved.ChatID)
	if err != nil {
		return nil, &vali.Varror{Error: ErrUnableToAddChatToUsersList}
	}

	return saved, nil
}
