package grpc_handlers

import (
	"context"

	"connectrpc.com/connect"
	grpc_helpers "github.com/kavkaco/Kavka-Core/delivery/grpc/helpers"
	grpc_model "github.com/kavkaco/Kavka-Core/delivery/grpc/model"
	"github.com/kavkaco/Kavka-Core/internal/service/message"
	"github.com/kavkaco/Kavka-Core/log"
	messagev1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/message"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/message/messagev1connect"
	"github.com/kavkaco/Kavka-Core/utils/vali"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/genproto/googleapis/rpc/code"
)

type MessageGrpcServer struct {
	logger         *log.SubLogger
	messageService message.MessageService
}

func NewMessageGrpcHandler(logger *log.SubLogger, messageService message.MessageService) messagev1connect.MessageServiceHandler {
	return MessageGrpcServer{logger, messageService}
}

func (h MessageGrpcServer) FetchMessages(ctx context.Context, req *connect.Request[messagev1.FetchMessagesRequest]) (*connect.Response[messagev1.FetchMessagesResponse], error) {
	chatID, err := primitive.ObjectIDFromHex(req.Msg.ChatId)
	if err != nil {
		return nil, grpc_helpers.GrpcVarror(&vali.Varror{Error: err}, connect.Code(code.Code_INTERNAL))
	}

	messages, varror := h.messageService.FetchMessages(ctx, chatID)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.Code(code.Code_INTERNAL))
	}

	res := connect.Response[messagev1.FetchMessagesResponse]{
		Msg: &messagev1.FetchMessagesResponse{
			Messages: grpc_model.TransformMessagesToGrpcModel(messages),
		},
	}

	return &res, nil
}
