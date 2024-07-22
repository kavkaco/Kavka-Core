package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"

	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/database"
	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	grpc_handlers "github.com/kavkaco/Kavka-Core/delivery/grpc/handlers"
	"github.com/kavkaco/Kavka-Core/delivery/grpc/interceptor"
	stream_consumers "github.com/kavkaco/Kavka-Core/infra/stream/consumers"
	stream_producers "github.com/kavkaco/Kavka-Core/infra/stream/producer"
	"github.com/kavkaco/Kavka-Core/internal/service/auth"
	"github.com/kavkaco/Kavka-Core/internal/service/chat"
	"github.com/kavkaco/Kavka-Core/pkg/email"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/auth/v1/authv1connect"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/chat/v1/chatv1connect"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/events/v1/eventsv1connect"
	"github.com/kavkaco/Kavka-Core/utils/hash"
	"github.com/rs/cors"
	auth_manager "github.com/tahadostifam/go-auth-manager"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleCORS(allowedOrigins []string, h http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:      allowedOrigins,
		AllowedMethods:      []string{"POST"},
		AllowedHeaders:      append(connectcors.AllowedHeaders(), []string{"X-Access-Token"}...),
		AllowPrivateNetwork: true,
	}).Handler(h)
}

func main() {
	ctx := context.Background()

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

	// Init Infra
	chatStreamEvents := make(chan map[string]interface{})
	chatProducer, err := stream_producers.NewBroadcastStreamProducer(&configs.Kafka)
	handleError(err)

	broadcastConsumer, err := stream_consumers.NewBroadcastConsumer(ctx, configs.Kafka.Brokers, *configs.Kafka.Sarama)
	handleError(err)

	hashManager := hash.NewHashManager(hash.DefaultHashParams)

	userRepo := repository_mongo.NewUserMongoRepository(mongoDB)
	// userService := user.NewUserService(userRepo)

	authRepo := repository_mongo.NewAuthMongoRepository(mongoDB)

	var emailService email.EmailService
	if config.CurrentEnv == config.Production {
		emailService = email.NewEmailService(&configs.Email, "email/templates")
	} else {
		emailService = email.NewEmailDevelopmentService()
	}

	authService := auth.NewAuthService(authRepo, userRepo, authManager, hashManager, emailService)

	chatRepo := repository_mongo.NewChatMongoRepository(mongoDB)
	chatService := chat.NewChatService(chatRepo, userRepo, chatProducer, chatStreamEvents)

	// messageRepo := repository.NewMessageRepository(mongoDB)
	// messageService := message.NewMessageService(messageRepo, chatRepo)

	// Init grpc server
	grpcListenAddr := fmt.Sprintf("%s:%d", configs.HTTP.Host, configs.HTTP.Port)
	gRPCRouter := http.NewServeMux()

	authInterceptor := interceptor.NewAuthInterceptor(authService)
	interceptors := connect.WithInterceptors(authInterceptor)

	authGrpcHandler := grpc_handlers.NewAuthGrpcHandler(authService)
	authGrpcRoute, authGrpcRouter := authv1connect.NewAuthServiceHandler(authGrpcHandler)

	chatGrpcHandler := grpc_handlers.NewChatGrpcHandler(chatService)
	chatGrpcRoute, chatGrpcRouter := chatv1connect.NewChatServiceHandler(chatGrpcHandler, interceptors)

	eventsGrpcHandler := grpc_handlers.NewEventsGrpcHandler(broadcastConsumer)
	eventsGrpcRoute, eventsGrpcRouter := eventsv1connect.NewEventsServiceHandler(eventsGrpcHandler, interceptors)

	gRPCRouter.Handle(authGrpcRoute, authGrpcRouter)
	gRPCRouter.Handle(chatGrpcRoute, chatGrpcRouter)
	gRPCRouter.Handle(eventsGrpcRoute, eventsGrpcRouter)

	// PPROF Memory Profiling Tool
	if config.CurrentEnv == config.Development {
		gRPCRouter.HandleFunc("/debug/pprof/*", pprof.Index)
		gRPCRouter.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	handler := handleCORS(configs.HTTP.Cors.AllowOrigins, gRPCRouter)
	server := &http.Server{
		Addr:         grpcListenAddr,
		Handler:      h2c.NewHandler(handler, &http2.Server{}),
		ReadTimeout:  0,
		WriteTimeout: 0,
		IdleTimeout:  0,
	}
	err = server.ListenAndServe()
	handleError(err)
}
