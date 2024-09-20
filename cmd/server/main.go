package main

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	connectcors "connectrpc.com/cors"
	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/database"
	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	"github.com/kavkaco/Kavka-Core/delivery/grpc"
	"github.com/kavkaco/Kavka-Core/infra/stream"
	"github.com/kavkaco/Kavka-Core/internal/service/auth"
	"github.com/kavkaco/Kavka-Core/internal/service/chat"
	"github.com/kavkaco/Kavka-Core/internal/service/message"
	"github.com/kavkaco/Kavka-Core/internal/service/search"
	"github.com/kavkaco/Kavka-Core/log"
	"github.com/kavkaco/Kavka-Core/pkg/email"

	"github.com/kavkaco/Kavka-Core/utils/hash"
	"github.com/rs/cors"
	auth_manager "github.com/tahadostifam/go-auth-manager"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func handleError(err error) {
	if err != nil {
		panic(err)
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
	// [=== Load Config ===]
	cfg := config.Read()

	// [=== Init Logger ===]
	log.InitGlobalLogger(cfg.Logger)

	// [=== Init MongoDB ===]
	mongoDB, err := database.GetMongoDBInstance(
		database.NewMongoDBConnectionString(
			cfg.Mongo.Host,
			cfg.Mongo.Port,
			cfg.Mongo.Username,
			cfg.Mongo.Password,
		),
		cfg.Mongo.DBName,
	)
	handleError(err)

	// [=== Init RedisDB ===]
	redisClient := database.GetRedisDBInstance(cfg.Redis)

	// [=== Init Auth Manager Service ===]
	authManager := auth_manager.NewAuthManager(redisClient, auth_manager.AuthManagerOpts{
		PrivateKey: cfg.Auth.SecretKey,
	})

	// [=== Init Infra ===]
	natsClient, err := stream.NewNATSAdapter(&cfg.Nats, log.NewSubLogger("infra"))
	handleError(err)

	streamPublisher, err := stream.NewStreamPublisher(natsClient)
	handleError(err)

	streamSubscriber, err := stream.NewStreamSubscriber(natsClient, log.NewSubLogger("stream-subscriber"))
	handleError(err)

	// [=== Init Internal Services & Repositories ===]
	hashManager := hash.NewHashManager(hash.DefaultHashParams)

	userRepo := repository_mongo.NewUserMongoRepository(mongoDB)
	// userService := user.NewUserService(userRepo)

	authRepo := repository_mongo.NewAuthMongoRepository(mongoDB)

	searchRepo := repository_mongo.NewSearchRepository(mongoDB)

	var emailService email.EmailService
	if config.CurrentEnv == config.Production {
		emailService = email.NewEmailService(&cfg.Email, "email/templates")
	} else {
		emailService = email.NewEmailDevelopmentService()
	}

	authService := auth.NewAuthService(authRepo, userRepo, authManager, hashManager, emailService)

	messageRepo := repository_mongo.NewMessageMongoRepository(mongoDB)

	chatRepo := repository_mongo.NewChatMongoRepository(mongoDB)
	chatService := chat.NewChatService(log.NewSubLogger("chat-service"), chatRepo, userRepo, messageRepo, streamPublisher)

	messageService := message.NewMessageService(log.NewSubLogger("message-service"), messageRepo, chatRepo, userRepo, streamPublisher)

	searchService := search.NewSearchService(log.NewSubLogger("search-service"), searchRepo)

	// [=== Init HTTP Server ===]
	router := http.NewServeMux()

	// [=== PPROF Memory Profiling Tool ===]
	if config.CurrentEnv == config.Development {
		router.HandleFunc("/debug/pprof/*", pprof.Index)
		router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	// [=== Init Grpc Server ===]
	grpc.NewGrpcServer(router, &grpc.Services{
		AuthService:      authService,
		ChatService:      chatService,
		MessageService:   messageService,
		SearchService:    searchService,
		StreamSubscriber: streamSubscriber,
	})

	// [=== Serve HTTP Server ===]
	handler := handleCORS(cfg.HTTP.Cors.AllowOrigins, router)
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:      h2c.NewHandler(handler, &http2.Server{}),
		ReadTimeout:  0,
		WriteTimeout: 0,
		IdleTimeout:  0,
	}
	err = server.ListenAndServe()
	handleError(err)
}
