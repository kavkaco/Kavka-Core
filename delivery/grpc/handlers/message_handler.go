package grpc_handlers

import (
	"context"

	"connectrpc.com/connect"
	grpc_helpers "github.com/kavkaco/Kavka-Core/delivery/grpc/helpers"
	"github.com/kavkaco/Kavka-Core/delivery/grpc/interceptor"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/model/proto_model_transformer"
	"github.com/kavkaco/Kavka-Core/internal/service/message"
	"github.com/kavkaco/Kavka-Core/log"

	"github.com/kavkaco/Kavka-Core/utils/vali"
	messagev1 "github.com/kavkaco/Kavka-ProtoBuf/gen/go/protobuf/message/v1"
	"github.com/kavkaco/Kavka-ProtoBuf/gen/go/protobuf/message/v1/messagev1connect"
	"google.golang.org/genproto/googleapis/rpc/code"
)

type MessageGrpcServer struct {
	logger         *log.SubLogger
	messageService *message.MessageService
}

func NewMessageGrpcHandler(logger *log.SubLogger, messageService *message.MessageService) messagev1connect.MessageServiceHandler {
	return MessageGrpcServer{logger, messageService}
}

func (h MessageGrpcServer) FetchMessages(ctx context.Context, req *connect.Request[messagev1.FetchMessagesRequest]) (*connect.Response[messagev1.FetchMessagesResponse], error) {
	chatID, err := model.ParseChatID(req.Msg.ChatId)
	if err != nil {
		return nil, grpc_helpers.GrpcVarror(&vali.Varror{Error: err}, connect.Code(code.Code_INTERNAL))
	}

	messages, varror := h.messageService.FetchMessages(ctx, chatID)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.Code(code.Code_INTERNAL))
	}

	res := connect.Response[messagev1.FetchMessagesResponse]{
		Msg: &messagev1.FetchMessagesResponse{
			Messages: proto_model_transformer.MessagesToProto(messages),
		},
	}

	return &res, nil
}

func (h MessageGrpcServer) SendTextMessage(ctx context.Context, req *connect.Request[messagev1.SendTextMessageRequest]) (*connect.Response[messagev1.SendTextMessageResponse], error) {
	userID := ctx.Value(interceptor.UserID{}).(model.UserID)
	if userID == "" {
		return nil, connect.NewError(connect.CodeDataLoss, interceptor.ErrEmptyUserID)
	}

	chatID, err := model.ParseChatID(req.Msg.ChatId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	msg, varror := h.messageService.SendTextMessage(ctx, chatID, userID, req.Msg.Text)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.CodePermissionDenied)
	}

	res := connect.Response[messagev1.SendTextMessageResponse]{
		Msg: &messagev1.SendTextMessageResponse{
			Message: proto_model_transformer.MessageToProto(msg),
		},
	}

	return &res, nil
}
