package main

import (
	"net/http"
	"net/http/pprof"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/database"
	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	"github.com/kavkaco/Kavka-Core/database/repo_sql"
	"github.com/kavkaco/Kavka-Core/delivery/grpc"
	"github.com/kavkaco/Kavka-Core/infra/cache"
	"github.com/kavkaco/Kavka-Core/infra/stream"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/internal/service/auth"
	"github.com/kavkaco/Kavka-Core/internal/service/chat"
	"github.com/kavkaco/Kavka-Core/internal/service/message"
	"github.com/kavkaco/Kavka-Core/internal/service/search"
	"github.com/kavkaco/Kavka-Core/internal/service/user"
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

	// [=== Init Infra ===]
	natsAdapter, err := stream.NewNATSAdapter(&cfg.Nats, log.NewSubLogger("infra"), stream.DefaultJetStreamConfig())
	handleError(err)
	defer natsAdapter.Close()

	streamPublisher, err := stream.NewStreamPublisher(natsAdapter)
	handleError(err)

	streamSubscriber, err := stream.NewStreamSubscriber(natsAdapter, log.NewSubLogger("stream-subscriber"))
	handleError(err)

	// [=== Init Database ===]
	var (
		userRepo    repository.UserRepository
		authRepo    repository.AuthRepository
		chatRepo    repository.ChatRepository
		messageRepo repository.MessageRepository
		searchRepo  repository.SearchRepository
	)

	switch cfg.SQL.Driver {
	case config.DriverMongo:
		log.Info("Database: MongoDB")

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

		userRepo = repository_mongo.NewUserMongoRepository(mongoDB)
		authRepo = repository_mongo.NewAuthMongoRepository(mongoDB)
		chatRepo = repository_mongo.NewChatMongoRepository(mongoDB)
		messageRepo = repository_mongo.NewMessageMongoRepository(mongoDB)
		searchRepo = repository_mongo.NewSearchRepository(mongoDB)

	case config.DriverPostgres, config.DriverSQLite:
		log.Info("Database: SQL")

		sqlAdapter, err := repo_sql.NewSQLAdapter(&cfg.SQL)
		handleError(err)
		defer sqlAdapter.Close()

		userRepo = repo_sql.NewSQLUserRepository(sqlAdapter.DB)
		authRepo = repo_sql.NewSQLAuthRepository(sqlAdapter.DB)
		chatRepo = repo_sql.NewSQLChatRepository(sqlAdapter.DB)
		messageRepo = repo_sql.NewSQLMessageRepository(sqlAdapter.DB)
		searchRepo = repo_sql.NewSQLSearchRepository(sqlAdapter.DB)

	default:
		panic(database.ErrUnsupportedDriver)
	}

	// [=== Init RedisDB ===]
	redisClient := database.GetRedisDBInstance(cfg.Redis)

	// [=== Init Auth Manager Service ===]
	authManager := auth_manager.NewAuthManager(redisClient, auth_manager.AuthManagerOpts{
		PrivateKey: cfg.Auth.SecretKey,
	})

	// [=== Init Cache ===]
	_ = cache.NewRedisCache(redisClient)

	// [=== Init Internal Services & Repositories ===]
	hashManager := hash.NewHashManager(hash.DefaultHashParams)

	_ = user.NewUserService(userRepo)

	var emailService email.EmailService
	if config.CurrentEnv == config.Production {
		emailService = email.NewEmailService(&cfg.Email, "email/templates")
	} else {
		emailService = email.NewEmailDevelopmentService()
	}

	authService := auth.NewAuthService(authRepo, userRepo, authManager, hashManager, emailService)

	chatService := chat.NewChatService(log.NewSubLogger("chat-service"), chatRepo, userRepo, messageRepo, streamPublisher)

	messageService := message.NewMessageService(
		log.NewSubLogger("message-service"),
		messageRepo,
		chatRepo,
		userRepo,
		streamPublisher,
		message.WithV2Schema(cfg.App.UseMessagesV2),
	)

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
