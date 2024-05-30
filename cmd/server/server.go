package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8089))
	handleError(err)

	grpcServer := grpc.NewServer()
	err = grpcServer.Serve(lis)
	handleError(err)
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// Init Zap Logger
	// logger := logs.InitZapLogger()

	// Load Configs
	// configs := config.Read()

	// Init MongoDB
	// mongoDB, mongoErr := database.GetMongoDBInstance(
	// 	database.NewMongoDBConnectionString(
	// 		configs.Mongo.Host,
	// 		configs.Mongo.Port,
	// 		configs.Mongo.Username,
	// 		configs.Mongo.Password,
	// 	),
	// 	configs.Mongo.DBName,
	// )
	// if mongoErr != nil {
	// 	panic(mongoErr)
	// }

	// Init RedisDB
	// redisClient := database.GetRedisDBInstance(configs.Redis)

	// authManager := auth_manager.NewAuthManager(redisClient, auth_manager.AuthManagerOpts{
	// 	PrivateKey: configs.App.Auth.SECRET,
	// })

	// hashManager := hash.NewHashManager(hash.DefaultHashParams)

	// Define paths
	// templatesPath := config.ProjectRootPath + "/app/views/mail/"
	// emailService := email.NewEmailService(logger, &configs.Email, templatesPath)

	// userRepo := repository.NewUserRepository(mongoDB)
	// userService := user.NewUserService(userRepo)

	// authRepo := repository.NewAuthRepository(mongoDB)
	// authService := auth.NewAuthService(authRepo, userRepo, authManager, hashManager)

	// chatRepo := repository.NewChatRepository(mongoDB)
	// chatService := chat.NewChatService(chatRepo, userRepo)

	// messageRepo := repository.NewMessageRepository(mongoDB)
	// messageService := message.NewMessageService(messageRepo, chatRepo)
}
