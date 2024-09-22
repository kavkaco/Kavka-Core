package main

import (
	"net/http"
	"net/http/pprof"

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
	auth_manager "github.com/tahadostifam/go-auth-manager"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
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
	err = grpc.NewGrpcServer(&cfg.HTTP, router, &grpc.Services{
		AuthService:      authService,
		ChatService:      chatService,
		MessageService:   messageService,
		SearchService:    searchService,
		StreamSubscriber: streamSubscriber,
	})
	handleError(err)
}
