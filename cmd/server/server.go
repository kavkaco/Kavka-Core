package main

import (
	"fmt"
	"log"
	"net"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/internal/service/auth"
	"github.com/kavkaco/Kavka-Core/pkg/auth_manager"
	auth_grpc "github.com/kavkaco/Kavka-Core/presentation/grpc/handlers"
	"github.com/kavkaco/Kavka-Core/presentation/grpc/pb"
	"github.com/kavkaco/Kavka-Core/utils/hash"
	"google.golang.org/grpc"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// Init Zap Logger
	// logger := logs.InitZapLogger()

	// Load Configs
	configs := config.Read()

	// Init MongoDB
	mongoDB, mongoErr := database.GetMongoDBInstance(
		database.NewMongoDBConnectionString(
			configs.Mongo.Host,
			configs.Mongo.Port,
			configs.Mongo.Username,
			configs.Mongo.Password,
		),
		configs.Mongo.DBName,
	)
	if mongoErr != nil {
		panic(mongoErr)
	}

	// Init RedisDB
	redisClient := database.GetRedisDBInstance(configs.Redis)

	authManager := auth_manager.NewAuthManager(redisClient, auth_manager.AuthManagerOpts{
		PrivateKey: configs.App.Auth.SECRET,
	})

	hashManager := hash.NewHashManager(hash.DefaultHashParams)

	// Define paths
	// templatesPath := config.ProjectRootPath + "/app/views/mail/"
	// emailService := email.NewEmailService(logger, &configs.Email, templatesPath)

	userRepo := repository.NewUserRepository(mongoDB)
	// userService := user.NewUserService(userRepo)

	authRepo := repository.NewAuthRepository(mongoDB)
	authService := auth.NewAuthService(authRepo, userRepo, authManager, hashManager)

	// chatRepo := repository.NewChatRepository(mongoDB)
	// chatService := chat.NewChatService(chatRepo, userRepo)

	// messageRepo := repository.NewMessageRepository(mongoDB)
	// messageService := message.NewMessageService(messageRepo, chatRepo)

	// Init gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8089))
	handleError(err)

	grpcServerRegistrar := grpc.NewServer()

	authGrpcServer := auth_grpc.NewAuthServerGrpc(grpcServerRegistrar, authService)
	pb.RegisterAuthServer(grpcServerRegistrar, &authGrpcServer)

	err = grpcServerRegistrar.Serve(lis)
	handleError(err)
}
