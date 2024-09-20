package grpc_handlers

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	grpc_helpers "github.com/kavkaco/Kavka-Core/delivery/grpc/helpers"
	"github.com/kavkaco/Kavka-Core/delivery/grpc/interceptor"

	"github.com/kavkaco/Kavka-Core/delivery/grpc/proto_model_transformer"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/service/chat"
	"github.com/kavkaco/Kavka-Core/log"
	chatv1 "github.com/kavkaco/Kavka-ProtoBuf/gen/go/protobuf/chat/v1"
	"github.com/kavkaco/Kavka-ProtoBuf/gen/go/protobuf/chat/v1/chatv1connect"
)

type chatHandler struct {
	logger      *log.SubLogger
	chatService *chat.ChatService
}

func NewChatGrpcHandler(logger *log.SubLogger, chatService *chat.ChatService) chatv1connect.ChatServiceHandler {
	return chatHandler{logger, chatService}
}

func (h chatHandler) GetDirectChat(ctx context.Context, req *connect.Request[chatv1.GetDirectChatRequest]) (*connect.Response[chatv1.GetDirectChatResponse], error) {
	userID := ctx.Value(interceptor.UserID{}).(model.UserID)
	if userID == "" {
		return nil, connect.NewError(connect.CodeDataLoss, interceptor.ErrEmptyUserID)
	}

	chatDto, varror := h.chatService.GetDirectChat(ctx, userID, req.Msg.RecipientUserId)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.CodeNotFound)
	}

	chatProto, err := proto_model_transformer.ChatToProto(*chatDto)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&chatv1.GetDirectChatResponse{
		Chat: chatProto,
	})

	return res, nil
}

func (h chatHandler) CreateChannel(ctx context.Context, req *connect.Request[chatv1.CreateChannelRequest]) (*connect.Response[chatv1.CreateChannelResponse], error) {
	userID := ctx.Value(interceptor.UserID{}).(model.UserID)
	if userID == "" {
		return nil, connect.NewError(connect.CodeDataLoss, interceptor.ErrEmptyUserID)
	}

	chat, varror := h.chatService.CreateChannel(ctx, userID, req.Msg.Title, req.Msg.Username, req.Msg.Description)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.CodeUnavailable)
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

func (h chatHandler) CreateDirect(ctx context.Context, req *connect.Request[chatv1.CreateDirectRequest]) (*connect.Response[chatv1.CreateDirectResponse], error) {
	userID := ctx.Value(interceptor.UserID{}).(model.UserID)
	if userID == "" {
		return nil, connect.NewError(connect.CodeDataLoss, interceptor.ErrEmptyUserID)
	}

	if req.Msg.RecipientUserId == userID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("user can't create direct chat with himself"))
	}

	chat, varror := h.chatService.CreateDirect(ctx, userID, req.Msg.RecipientUserId)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.CodeUnavailable)
	}

	chatProto, err := proto_model_transformer.ChatToProto(*chat)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&chatv1.CreateDirectResponse{
		Chat: chatProto,
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
		return nil, grpc_helpers.GrpcVarror(varror, connect.CodeUnavailable)
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
	userID := ctx.Value(interceptor.UserID{}).(model.UserID)
	if userID == "" {
		return nil, connect.NewError(connect.CodeDataLoss, interceptor.ErrEmptyUserID)
	}

	chatID, err := model.ParseChatID(req.Msg.ChatId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	chat, varror := h.chatService.GetChat(ctx, userID, chatID)
	if varror != nil {
		return nil, varror.Error
	}

	chatGetter := model.NewChatDTO(chat)
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

func (h chatHandler) JoinChat(ctx context.Context, req *connect.Request[chatv1.JoinChatRequest]) (*connect.Response[chatv1.JoinChatResponse], error) {
	userID := ctx.Value(interceptor.UserID{}).(model.UserID)
	if userID == "" {
		return nil, connect.NewError(connect.CodeDataLoss, interceptor.ErrEmptyUserID)
	}

	chatID, err := model.ParseChatID(req.Msg.ChatId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	joinResult, varror := h.chatService.JoinChat(ctx, chatID, userID)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.CodeInternal)
	}

	if !joinResult.Joined {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("joining chat failed"))
	}

	protoChat, err := proto_model_transformer.ChatToProto(*joinResult.UpdatedChat)
	if err != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.CodeInternal)
	}

	res := &connect.Response[chatv1.JoinChatResponse]{Msg: &chatv1.JoinChatResponse{
		Chat: protoChat,
	}}

	return res, nil
}
