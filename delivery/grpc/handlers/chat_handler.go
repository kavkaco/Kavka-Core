package grpc_handlers

import (
	"context"

	"connectrpc.com/connect"
	grpc_helpers "github.com/kavkaco/Kavka-Core/delivery/grpc/helpers"
	"github.com/kavkaco/Kavka-Core/delivery/grpc/interceptor"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/service/chat"
	"github.com/kavkaco/Kavka-Core/log"
	chatv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/chat/v1"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/chat/v1/chatv1connect"
	"github.com/kavkaco/Kavka-Core/protobuf/proto_model_transformer"
)

type chatHandler struct {
	logger      *log.SubLogger
	chatService chat.ChatService
}

func NewChatGrpcHandler(logger *log.SubLogger, chatService chat.ChatService) chatv1connect.ChatServiceHandler {
	return chatHandler{logger, chatService}
}

func (h chatHandler) CreateChannel(ctx context.Context, req *connect.Request[chatv1.CreateChannelRequest]) (*connect.Response[chatv1.CreateChannelResponse], error) {
	userID := ctx.Value(interceptor.UserID{}).(model.UserID)
	if userID == "" {
		return nil, connect.NewError(connect.CodeDataLoss, interceptor.ErrEmptyUserID)
	}

	chat, varror := h.chatService.CreateChannel(ctx, userID, req.Msg.Title, req.Msg.Username, req.Msg.Description)
	if varror != nil {
		connectErr := connect.NewError(connect.CodeUnavailable, varror.Error)
		varrorDetail, _ := grpc_helpers.VarrorAsGrpcErrDetails(varror)
		connectErr.AddDetail(varrorDetail)
		return nil, connectErr
	}

	chatProto, err := proto_model_transformer.ChatToProto(*chat)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&chatv1.CreateChannelResponse{
		Chat: chatProto,
	})

	return res, nil
}

// FIXME - Later we will work on this, no problem at the moment!
//
// I think we must call 2 diff methods for peer to peer messaging...
// user first creates the direct chat and then can send messages,
// or with custom web clients or etc they even can create a direct chat with no any messages included...
func (h chatHandler) CreateDirect(ctx context.Context, req *connect.Request[chatv1.CreateDirectRequest]) (*connect.Response[chatv1.CreateDirectResponse], error) {
	res := connect.NewResponse(&chatv1.CreateDirectResponse{
		Chat: nil,
	})
	return res, nil
}

func (h chatHandler) CreateGroup(ctx context.Context, req *connect.Request[chatv1.CreateGroupRequest]) (*connect.Response[chatv1.CreateGroupResponse], error) {
	userID := ctx.Value(interceptor.UserID{}).(model.UserID)
	if userID == "" {
		return nil, connect.NewError(connect.CodeDataLoss, interceptor.ErrEmptyUserID)
	}

	chat, varror := h.chatService.CreateGroup(ctx, userID, req.Msg.Title, req.Msg.Username, req.Msg.Description)
	if varror != nil {
		connectErr := connect.NewError(connect.CodeUnavailable, varror.Error)
		varrorDetail, _ := grpc_helpers.VarrorAsGrpcErrDetails(varror)
		connectErr.AddDetail(varrorDetail)
		return nil, connectErr
	}

	chatProto, err := proto_model_transformer.ChatToProto(*chat)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&chatv1.CreateGroupResponse{
		Chat: chatProto,
	})

	return res, nil
}

func (h chatHandler) GetChat(ctx context.Context, req *connect.Request[chatv1.GetChatRequest]) (*connect.Response[chatv1.GetChatResponse], error) {
	chatID, err := model.ParseChatID(req.Msg.ChatId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	chat, varror := h.chatService.GetChat(ctx, chatID)
	if varror != nil {
		return nil, varror.Error
	}

	chatGetter := model.NewChatGetter(chat)
	chatProto, err := proto_model_transformer.ChatToProto(*chatGetter)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := &connect.Response[chatv1.GetChatResponse]{
		Msg: &chatv1.GetChatResponse{
			Chat: chatProto,
		},
	}

	return res, nil
}

func (h chatHandler) GetUserChats(ctx context.Context, req *connect.Request[chatv1.GetUserChatsRequest]) (*connect.Response[chatv1.GetUserChatsResponse], error) {
	userID := ctx.Value(interceptor.UserID{}).(model.UserID)
	if userID == "" {
		return nil, connect.NewError(connect.CodeDataLoss, interceptor.ErrEmptyUserID)
	}

	chats, varror := h.chatService.GetUserChats(ctx, userID)
	if varror != nil {
		return nil, varror.Error
	}

	chatsProto, err := proto_model_transformer.ChatsToProto(chats)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(
		&chatv1.GetUserChatsResponse{
			Chats: chatsProto,
		},
	)

	return res, nil
}
