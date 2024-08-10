package grpc_handlers

import (
	"context"

	"connectrpc.com/connect"
	grpc_helpers "github.com/kavkaco/Kavka-Core/delivery/grpc/helpers"

	"github.com/kavkaco/Kavka-Core/internal/service/auth"
	authv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/auth/v1"
	"github.com/kavkaco/Kavka-Core/protobuf/proto_model_transformer"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/protobuf/types/known/durationpb"
)

type AuthGrpcServer struct {
	authService auth.AuthService
}

func NewAuthGrpcHandler(authService auth.AuthService) AuthGrpcServer {
	return AuthGrpcServer{authService}
}

func (a AuthGrpcServer) Login(ctx context.Context, req *connect.Request[authv1.LoginRequest]) (*connect.Response[authv1.LoginResponse], error) {
	user, accessToken, refreshToken, varror := a.authService.Login(ctx, req.Msg.Email, req.Msg.Password)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.Code(code.Code_INTERNAL))
	}

	res := connect.NewResponse(&authv1.LoginResponse{
		User:         proto_model_transformer.UserToProto(*user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	return res, nil
}

func (a AuthGrpcServer) Register(ctx context.Context, req *connect.Request[authv1.RegisterRequest]) (*connect.Response[authv1.RegisterResponse], error) {
	_, varror := a.authService.Register(ctx, req.Msg.Name, req.Msg.LastName, req.Msg.Username, req.Msg.Email, req.Msg.Password, req.Msg.VerifyEmailRedirectUrl)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.Code(code.Code_INVALID_ARGUMENT))
	}

	res := connect.NewResponse(&authv1.RegisterResponse{})

	return res, nil
}

func (a AuthGrpcServer) Authenticate(ctx context.Context, req *connect.Request[authv1.AuthenticateRequest]) (*connect.Response[authv1.AuthenticateResponse], error) {
	user, varror := a.authService.Authenticate(ctx, req.Msg.AccessToken)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.Code(code.Code_PERMISSION_DENIED))
	}

	res := connect.NewResponse(&authv1.AuthenticateResponse{
		User: proto_model_transformer.UserToProto(*user),
	})

	return res, nil
}

func (a AuthGrpcServer) ChangePassword(ctx context.Context, req *connect.Request[authv1.ChangePasswordRequest]) (*connect.Response[authv1.ChangePasswordResponse], error) {
	varror := a.authService.ChangePassword(ctx, req.Msg.AccessToken, req.Msg.OldPassword, req.Msg.NewPassword)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.CodeUnavailable)
	}

	res := connect.NewResponse(&authv1.ChangePasswordResponse{})
	return res, nil
}

func (a AuthGrpcServer) RefreshToken(ctx context.Context, req *connect.Request[authv1.RefreshTokenRequest]) (*connect.Response[authv1.RefreshTokenResponse], error) {
	newAccessToken, varror := a.authService.RefreshToken(ctx, req.Msg.UserId, req.Msg.RefreshToken)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.Code(code.Code_UNAVAILABLE))
	}

	res := connect.NewResponse(&authv1.RefreshTokenResponse{
		AccessToken: newAccessToken,
	})

	return res, nil
}

func (a AuthGrpcServer) SendResetPassword(ctx context.Context, req *connect.Request[authv1.SendResetPasswordRequest]) (*connect.Response[authv1.SendResetPasswordResponse], error) {
	resetPasswordToken, timeout, varror := a.authService.SendResetPassword(ctx, req.Msg.Email, req.Msg.ResetPasswordRedirectUrl)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.Code(code.Code_UNAVAILABLE))
	}

	timeoutProto := durationpb.New(timeout)
	res := connect.NewResponse(&authv1.SendResetPasswordResponse{
		ResetPasswordToken: resetPasswordToken,
		Timeout:            timeoutProto,
	})

	return res, nil
}

func (a AuthGrpcServer) SubmitResetPassword(ctx context.Context, req *connect.Request[authv1.SubmitResetPasswordRequest]) (*connect.Response[authv1.SubmitResetPasswordResponse], error) {
	varror := a.authService.SubmitResetPassword(ctx, req.Msg.ResetPasswordToken, req.Msg.NewPassword)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.Code(code.Code_UNAVAILABLE))
	}

	res := connect.NewResponse(&authv1.SubmitResetPasswordResponse{})
	return res, nil
}

func (a AuthGrpcServer) VerifyEmail(ctx context.Context, req *connect.Request[authv1.VerifyEmailRequest]) (*connect.Response[authv1.VerifyEmailResponse], error) {
	varror := a.authService.VerifyEmail(ctx, req.Msg.VerifyEmailToken)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.Code(code.Code_UNAVAILABLE))
	}

	res := connect.NewResponse(&authv1.VerifyEmailResponse{})

	return res, nil
}
