package grpc_service

import (
	"context"

	"connectrpc.com/connect"
	grpc_common "github.com/kavkaco/Kavka-Core/delivery/grpc/common"
	"github.com/kavkaco/Kavka-Core/internal/service/auth"
	authv1 "github.com/kavkaco/Kavka-ProtoBuf/gen/go/proto/auth/v1"
)

type AuthGrpcServer struct {
	authService auth.AuthService
}

func NewAuthGrpcHandler(authService auth.AuthService) AuthGrpcServer {
	return AuthGrpcServer{authService}
}

// Login implements authv1connect.AuthServiceHandler.
func (a AuthGrpcServer) Login(ctx context.Context, req *connect.Request[authv1.LoginRequest]) (*connect.Response[authv1.LoginResponse], error) {
	user, accessToken, refreshToken, err := a.authService.Login(ctx, req.Msg.Email, req.Msg.Password)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, err)
	}

	res := connect.NewResponse(&authv1.LoginResponse{
		User:         grpc_common.TransformUserToGrpcModel(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	return res, nil
}

// Authenticate implements authv1connect.AuthServiceHandler.
func (a AuthGrpcServer) Authenticate(context.Context, *connect.Request[authv1.AuthenticateRequest]) (*connect.Response[authv1.AuthenticateResponse], error) {
	panic("unimplemented")
}

// ChangePassword implements authv1connect.AuthServiceHandler.
func (a AuthGrpcServer) ChangePassword(context.Context, *connect.Request[authv1.ChangePasswordRequest]) (*connect.Response[authv1.ChangePasswordResponse], error) {
	panic("unimplemented")
}

// RefreshToken implements authv1connect.AuthServiceHandler.
func (a AuthGrpcServer) RefreshToken(context.Context, *connect.Request[authv1.RefreshTokenRequest]) (*connect.Response[authv1.RefreshTokenResponse], error) {
	panic("unimplemented")
}

// Register implements authv1connect.AuthServiceHandler.
func (a AuthGrpcServer) Register(context.Context, *connect.Request[authv1.RegisterRequest]) (*connect.Response[authv1.RegisterResponse], error) {
	panic("unimplemented")
}

// SendResetPasswordVerification implements authv1connect.AuthServiceHandler.
func (a AuthGrpcServer) SendResetPasswordVerification(context.Context, *connect.Request[authv1.SendResetPasswordVerificationRequest]) (*connect.Response[authv1.SendResetPasswordVerificationResponse], error) {
	panic("unimplemented")
}

// SubmitResetPassword implements authv1connect.AuthServiceHandler.
func (a AuthGrpcServer) SubmitResetPassword(context.Context, *connect.Request[authv1.SubmitResetPasswordRequest]) (*connect.Response[authv1.SubmitResetPasswordResponse], error) {
	panic("unimplemented")
}

// VerifyEmail implements authv1connect.AuthServiceHandler.
func (a AuthGrpcServer) VerifyEmail(context.Context, *connect.Request[authv1.VerifyEmailRequest]) (*connect.Response[authv1.VerifyEmailResponse], error) {
	panic("unimplemented")
}
