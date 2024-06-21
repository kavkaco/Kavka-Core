package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/database"
	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	grpc_service "github.com/kavkaco/Kavka-Core/delivery/grpc"
	"github.com/kavkaco/Kavka-Core/internal/service/auth"
	"github.com/kavkaco/Kavka-Core/pkg/auth_manager"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/auth/v1/authv1connect"
	"github.com/kavkaco/Kavka-Core/utils/hash"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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
	mux := http.NewServeMux()

	authGrpcHandler := grpc_service.NewAuthGrpcHandler(authService)
	path, handler := authv1connect.NewAuthServiceHandler(authGrpcHandler)

	mux.Handle(path, handler)

	server := &http.Server{
		Addr:         grpcListenAddr,
		Handler:      h2c.NewHandler(mux, &http2.Server{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 2 * time.Second,
	}
	err := server.ListenAndServe()
	handleError(err)
}
