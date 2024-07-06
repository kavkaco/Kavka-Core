package grpc_service

import (
	"context"

	"connectrpc.com/connect"
	grpc_model "github.com/kavkaco/Kavka-Core/delivery/grpc/model"
	"github.com/kavkaco/Kavka-Core/internal/service/auth"
	authv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/auth/v1"
	"google.golang.org/protobuf/types/known/durationpb"
)

type AuthGrpcServer struct {
	authService auth.AuthService
}

func NewAuthGrpcHandler(authService auth.AuthService) AuthGrpcServer {
	return AuthGrpcServer{authService}
}

func (a AuthGrpcServer) Login(ctx context.Context, req *connect.Request[authv1.LoginRequest]) (*connect.Response[authv1.LoginResponse], error) {
	user, accessToken, refreshToken, err := a.authService.Login(ctx, req.Msg.Email, req.Msg.Password)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, err)
	}

	res := connect.NewResponse(&authv1.LoginResponse{
		User:         grpc_model.TransformUserToGrpcModel(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	return res, nil
}

func (a AuthGrpcServer) Register(ctx context.Context, req *connect.Request[authv1.RegisterRequest]) (*connect.Response[authv1.RegisterResponse], error) {
	_, _, err := a.authService.Register(ctx, req.Msg.Name, req.Msg.LastName, req.Msg.Username, req.Msg.Email, req.Msg.Password, req.Msg.VerifyEmailRedirectUrl)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := connect.NewResponse(&authv1.RegisterResponse{})

	return res, nil
}

func (a AuthGrpcServer) Authenticate(ctx context.Context, req *connect.Request[authv1.AuthenticateRequest]) (*connect.Response[authv1.AuthenticateResponse], error) {
	user, err := a.authService.Authenticate(ctx, req.Msg.AccessToken)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, err)
	}

	res := connect.NewResponse(&authv1.AuthenticateResponse{
		User: grpc_model.TransformUserToGrpcModel(user),
	})

	return res, nil
}

func (a AuthGrpcServer) ChangePassword(ctx context.Context, req *connect.Request[authv1.ChangePasswordRequest]) (*connect.Response[authv1.ChangePasswordResponse], error) {
	err := a.authService.ChangePassword(ctx, req.Msg.AccessToken, req.Msg.OldPassword, req.Msg.NewPassword)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := connect.NewResponse(&authv1.ChangePasswordResponse{})
	return res, nil
}

func (a AuthGrpcServer) RefreshToken(ctx context.Context, req *connect.Request[authv1.RefreshTokenRequest]) (*connect.Response[authv1.RefreshTokenResponse], error) {
	newAccessToken, err := a.authService.RefreshToken(ctx, req.Msg.RefreshToken, req.Msg.AccessToken)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := connect.NewResponse(&authv1.RefreshTokenResponse{
		AccessToken: newAccessToken,
	})

	return res, nil
}

func (a AuthGrpcServer) SendResetPassword(ctx context.Context, req *connect.Request[authv1.SendResetPasswordRequest]) (*connect.Response[authv1.SendResetPasswordResponse], error) {
	resetPasswordToken, timeout, err := a.authService.SendResetPassword(ctx, req.Msg.Email, req.Msg.ResetPasswordRedirectUrl)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	timeoutProto := durationpb.New(timeout)
	res := connect.NewResponse(&authv1.SendResetPasswordResponse{
		ResetPasswordToken: resetPasswordToken,
		Timeout:            timeoutProto,
	})

	return res, nil
}

func (a AuthGrpcServer) SubmitResetPassword(ctx context.Context, req *connect.Request[authv1.SubmitResetPasswordRequest]) (*connect.Response[authv1.SubmitResetPasswordResponse], error) {
	err := a.authService.SubmitResetPassword(ctx, req.Msg.ResetPasswordToken, req.Msg.NewPassword)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := connect.NewResponse(&authv1.SubmitResetPasswordResponse{})
	return res, nil
}

func (a AuthGrpcServer) VerifyEmail(ctx context.Context, req *connect.Request[authv1.VerifyEmailRequest]) (*connect.Response[authv1.VerifyEmailResponse], error) {
	err := a.authService.VerifyEmail(ctx, req.Msg.VerifyEmailToken)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&authv1.VerifyEmailResponse{})

	return res, nil
}
