package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/database"
	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	grpc_service "github.com/kavkaco/Kavka-Core/delivery/grpc"
	"github.com/kavkaco/Kavka-Core/internal/service/auth"
	"github.com/kavkaco/Kavka-Core/pkg/auth_manager"
	"github.com/kavkaco/Kavka-Core/utils/hash"
	"github.com/kavkaco/Kavka-ProtoBuf/gen/go/proto/auth/v1/authv1connect"
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
		PrivateKey: configs.Auth.SecretKey,
	})

	hashManager := hash.NewHashManager(hash.DefaultHashParams)

	// Define paths
	// templatesPath := config.ProjectRootPath + "/app/views/mail/"
	// emailService := email.NewEmailService(logger, &configs.Email, templatesPath)

	userRepo := repository_mongo.NewUserMongoRepository(mongoDB)
	// userService := user.NewUserService(userRepo)

	authRepo := repository_mongo.NewAuthMongoRepository(mongoDB)
	authService := auth.NewAuthService(authRepo, userRepo, authManager, hashManager)

	// chatRepo := repository.NewChatRepository(mongoDB)
	// chatService := chat.NewChatService(chatRepo, userRepo)

	// messageRepo := repository.NewMessageRepository(mongoDB)
	// messageService := message.NewMessageService(messageRepo, chatRepo)

	// Init grpc server

	grpcListenAddr := fmt.Sprintf("%s:%d", configs.HTTP.Host, configs.HTTP.Port)
	app := fiber.New()

	authGrpcHandler := grpc_service.NewAuthGrpcHandler(authService)
	authPath, authHandler := authv1connect.NewAuthServiceHandler(authGrpcHandler)

	app.All(authPath, adaptor.HTTPHandler(authHandler))

	err := app.Listen(grpcListenAddr)
	handleError(err)

	// grpcListenAddr := fmt.Sprintf("%s:%d", configs.HTTP.Host, configs.HTTP.Port)
	// mux := http.NewServeMux()

	// authGrpcHandler := grpc_service.NewAuthGrpcHandler(authService)
	// path, handler := authv1connect.NewAuthServiceHandler(authGrpcHandler)

	// mux.Handle(path, handler)

	// err := http.ListenAndServe(grpcListenAddr, h2c.NewHandler(mux, &http2.Server{}))
	// handleError(err)
}
