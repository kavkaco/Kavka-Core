package grpc_service

import (
	"context"

	"connectrpc.com/connect"
	"github.com/kavkaco/Kavka-Core/delivery/grpc/interceptor"
	grpc_model "github.com/kavkaco/Kavka-Core/delivery/grpc/model"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/service/chat"
	chatv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/chat/v1"
	v1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/model/chat/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/chat/v1/chatv1connect"
)

type handler struct {
	chatService chat.ChatService
}

func NewChatGrpcHandler(chatService chat.ChatService) chatv1connect.ChatServiceHandler {
	return handler{chatService}
}

func (h handler) CreateChannel(ctx context.Context, req *connect.Request[chatv1.CreateChannelRequest]) (*connect.Response[chatv1.CreateChannelResponse], error) {
	userID := ctx.Value(interceptor.UserIDKey{}).(model.UserID)

	chat, err := h.chatService.CreateChannel(ctx, userID, req.Msg.Title, req.Msg.Username, req.Msg.Description)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	chatGrpcModel, err := grpc_model.TransformChatToGrpcModel(chat)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&chatv1.CreateChannelResponse{
		Chat: chatGrpcModel,
	})

	return res, nil
}

func (h handler) CreateDirect(ctx context.Context, req *connect.Request[chatv1.CreateDirectRequest]) (*connect.Response[chatv1.CreateDirectResponse], error) {
	chat, err := h.chatService.CreateDirect(ctx, req.Msg.UserId, req.Msg.RecipientUserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	chatGrpcModel, err := grpc_model.TransformChatToGrpcModel(chat)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	res := connect.NewResponse(&chatv1.CreateDirectResponse{
		Chat: chatGrpcModel,
	})
	return res, nil
}

func (h handler) CreateGroup(ctx context.Context, req *connect.Request[chatv1.CreateGroupRequest]) (*connect.Response[chatv1.CreateGroupResponse], error) {
	chat, err := h.chatService.CreateGroup(ctx, req.Msg.UserId, req.Msg.Title, req.Msg.Title, req.Msg.Description)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	chatGrpcModel, err := grpc_model.TransformChatToGrpcModel(chat)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	res := connect.NewResponse(&chatv1.CreateGroupResponse{
		Chat: chatGrpcModel,
	})
	return res, nil
}

func (h handler) GetChat(ctx context.Context, req *connect.Request[chatv1.GetChatRequest]) (*connect.Response[chatv1.GetChatResponse], error) {
	chatID, err := primitive.ObjectIDFromHex(req.Msg.ChatId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	chat, err := h.chatService.GetChat(ctx, chatID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	chatGrpcModel, err := grpc_model.TransformChatToGrpcModel(chat)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	res := connect.NewResponse(&chatv1.GetChatResponse{
		Chat: chatGrpcModel,
	})
	return res, nil
}

func (h handler) GetUserChats(ctx context.Context, req *connect.Request[chatv1.GetUserChatsRequest]) (*connect.Response[chatv1.GetUserChatsResponse], error) {
	chats, err := h.chatService.GetUserChats(ctx, req.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	var chatsGrpcModel []*v1.Chat

	for _, v := range chats {
		chatGrpcModel, err := grpc_model.TransformChatToGrpcModel(&v)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		chatsGrpcModel = append(chatsGrpcModel, chatGrpcModel)
	}
	res := connect.NewResponse(&chatv1.GetUserChatsResponse{
		Chats: chatsGrpcModel,
	})
	return res, nil
}
