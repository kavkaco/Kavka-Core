package grpc_handlers

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/kavkaco/Kavka-Core/delivery/grpc/interceptor"
	"github.com/kavkaco/Kavka-Core/internal/service/user"
	"github.com/kavkaco/Kavka-Core/log"
	userv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/user/v1"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/user/v1/userv1connect"
)

type userGrpcServer struct {
	logger      *log.SubLogger
	userService user.UserService
}

func NewUserGrpcServer(logger *log.SubLogger, userService user.UserService) userv1connect.UserServiceHandler {
	return &userGrpcServer{logger, userService}
}

func (h *userGrpcServer) UploadProfile(ctx context.Context, stream *connect.ClientStream[userv1.UploadProfileRequest]) (*connect.Response[userv1.UploadProfileResponse], error) {
	headers := stream.RequestHeader()

	filename := headers.Get("filename")
	if filename == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing filename"))
	}
	userID := ctx.Value(interceptor.UserID{}).(string)
	if userID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing userID"))
	}

	for stream.Receive() {
		bytes := stream.Msg().Bytes

		err := h.userService.UpdateProfilePicture(ctx, userID, filename, bytes)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err.Error)
		}
	}

	res := connect.NewResponse(&userv1.UploadProfileResponse{})

	return res, nil
}
