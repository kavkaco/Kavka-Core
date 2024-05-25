package main

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kavkaco/Kavka-Core/app/middleware"
	"github.com/kavkaco/Kavka-Core/app/router"
	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/internal/service/auth"
	"github.com/kavkaco/Kavka-Core/internal/service/chat"
	"github.com/kavkaco/Kavka-Core/internal/service/message"
	user "github.com/kavkaco/Kavka-Core/internal/service/user"
	"github.com/kavkaco/Kavka-Core/logs"
	"github.com/kavkaco/Kavka-Core/pkg/auth_manager"
	"github.com/kavkaco/Kavka-Core/pkg/email"
	"github.com/kavkaco/Kavka-Core/socket/adapters"
	"github.com/kavkaco/Kavka-Core/socket/handlers"
	"github.com/kavkaco/Kavka-Core/utils/hash"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init Zap Logger
	logger := logs.InitZapLogger()

	// Define paths
	TemplatesPath := config.ProjectRootPath + "/app/views/mail/"

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

	// Init WebServer
	ginEngine := gin.New()
	api := ginEngine.Group("/api/v1")

	// Cors
	api.Use(cors.New(cors.Config{
		AllowOrigins:     []string{configs.App.Server.CORS.AllowOrigins},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Refresh", "Authorization"},
		ExposeHeaders:    []string{"Refresh", "Authorization"},
		AllowCredentials: true,
	}))

	authManager := auth_manager.NewAuthManager(redisClient, auth_manager.AuthManagerOpts{
		PrivateKey: configs.App.Auth.SECRET,
	})

	hashManager := hash.NewHashManager(hash.DefaultHashParams)

	emailService := email.NewEmailService(logger, &configs.Email, TemplatesPath)

	userRepo := repository.NewUserRepository(mongoDB)
	userService := user.NewUserService(userRepo)

	authRepo := repository.NewAuthRepository(mongoDB)
	authService := auth.NewAuthService(authRepo, userRepo, authManager, hashManager)

	chatRepo := repository.NewChatRepository(mongoDB)
	chatService := chat.NewChatService(chatRepo, userRepo)

	messageRepo := repository.NewMessageRepository(mongoDB)
	messageService := message.NewMessageService(messageRepo, chatRepo)

	// Init routes
	router.NewAuthRouter(ctx, logger, api.Group("/auth"), authService, chatService, *emailService)

	// Init websocket server
	websocketAdapter := adapters.NewWebsocketAdapter(logger)

	handlerServices := handlers.HandlerServices{
		UserService:    userService,
		ChatService:    chatService,
		MessageService: messageService,
	}

	api.GET("/ws", middleware.AuthenticatedMiddleware(ctx, authService), router.WebsocketRoute(ctx, logger, websocketAdapter, handlerServices))

	// Everything almost done!
	err := ginEngine.Run(configs.App.HTTP.Address)
	if err != nil {
		panic(err)
	}
}
