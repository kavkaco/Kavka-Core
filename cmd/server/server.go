package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/database"
	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	grpc_service "github.com/kavkaco/Kavka-Core/delivery/grpc"
	"github.com/kavkaco/Kavka-Core/delivery/grpc/interceptor"
	"github.com/kavkaco/Kavka-Core/internal/service/auth"
	"github.com/kavkaco/Kavka-Core/internal/service/chat"
	"github.com/kavkaco/Kavka-Core/pkg/auth_manager"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/auth/v1/authv1connect"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/chat/v1/chatv1connect"
	"github.com/kavkaco/Kavka-Core/utils/hash"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func newRouterCORS(currentEnv config.Env, prodCORS []string, h http.Handler) http.Handler {
	opts := cors.Options{
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	}

	if currentEnv == config.Production {
		opts.AllowedOrigins = prodCORS
	} else {
		opts.AllowedOrigins = []string{"*"}
	}

	c := cors.New(opts)

	return c.Handler(h)
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

	chatRepo := repository_mongo.NewChatMongoRepository(mongoDB)
	chatService := chat.NewChatService(chatRepo, userRepo)

	// messageRepo := repository.NewMessageRepository(mongoDB)
	// messageService := message.NewMessageService(messageRepo, chatRepo)

	// Init grpc server
	grpcListenAddr := fmt.Sprintf("%s:%d", configs.HTTP.Host, configs.HTTP.Port)
	gRPCRouter := http.NewServeMux()

	authInterceptor := interceptor.NewAuthInterceptor(authService)
	interceptors := connect.WithInterceptors(authInterceptor)

	authGrpcHandler := grpc_service.NewAuthGrpcHandler(authService)
	authGrpcRoute, authGrpcRouter := authv1connect.NewAuthServiceHandler(authGrpcHandler)

	chatGrpcHandler := grpc_service.NewChatGrpcHandler(chatService)
	chatGrpcRoute, chatGrpcRouter := chatv1connect.NewChatServiceHandler(chatGrpcHandler, interceptors)

	gRPCRouter.Handle(authGrpcRoute, authGrpcRouter)
	gRPCRouter.Handle(chatGrpcRoute, chatGrpcRouter)

	configCors := []string{"*"} // FIXME - after refactoring config pkg
	gRPCHandler := newRouterCORS(config.CurrentEnv, configCors, gRPCRouter)
	server := &http.Server{
		Addr:         grpcListenAddr,
		Handler:      h2c.NewHandler(gRPCHandler, &http2.Server{}),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 7 * time.Second,
	}
	err := server.ListenAndServe()
	handleError(err)
}
