package auth_grpc

import (
	"context"

	"github.com/kavkaco/Kavka-Core/delivery/grpc/pb"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGrpcServer struct {
	pb.UnimplementedAuthServer
	authService auth.AuthService
}

func NewAuthServerGrpc(gs *grpc.Server, authService auth.AuthService) AuthGrpcServer {
	return AuthGrpcServer{authService: authService}
}

func transformUserToGrpc(user *model.User) *pb.User {
	return &pb.User{
		UserId:    user.UserID,
		Name:      user.Name,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
		Biography: user.Biography,
	}
}

// func transformGrpcToUser(user *pb.User) *model.User {
// 	return &model.User{
// 		UserID:    user.UserId,
// 		Name:      user.Name,
// 		LastName:  user.LastName,
// 		Email:     user.Email,
// 		Username:  user.Username,
// 		Biography: user.Biography,
// 	}
// }

func (s *AuthGrpcServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, accessToken, refreshToken, err := s.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	return &pb.LoginResponse{
		User:         transformUserToGrpc(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
