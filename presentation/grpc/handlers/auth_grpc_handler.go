package auth_grpc

import (
	"context"
	"fmt"

	"github.com/kavkaco/Kavka-Core/internal/service/auth"
	"github.com/kavkaco/Kavka-Core/presentation/grpc/pb"
	"google.golang.org/grpc"
)

type AuthGrpcServer struct {
	authService auth.AuthService
	pb.UnimplementedAuthHandlerServer
}

func NewAuthServerGrpc(gs *grpc.Server, authService auth.AuthService) *AuthGrpcServer {
	return &AuthGrpcServer{}
}

func (s *AuthGrpcServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	fmt.Println("Request email: " + req.Email)

	return &pb.LoginResponse{
		User:         nil,
		AccessToken:  "access",
		RefreshToken: "refresh",
	}, nil
}
