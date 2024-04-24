package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kavkaco/Kavka-Core/app/router"
	"github.com/kavkaco/Kavka-Core/app/socket"
	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/database"
	chatRepository "github.com/kavkaco/Kavka-Core/internal/repository/chat"
	messageRepository "github.com/kavkaco/Kavka-Core/internal/repository/message"
	userRepository "github.com/kavkaco/Kavka-Core/internal/repository/user"
	"github.com/kavkaco/Kavka-Core/internal/service"
	"github.com/kavkaco/Kavka-Core/logs"
	"github.com/kavkaco/Kavka-Core/pkg/session"
	"github.com/kavkaco/Kavka-Core/pkg/sms_service"
)

func main() {
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
	app := gin.New()

	// Cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{configs.App.Server.CORS.AllowOrigins},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Refresh", "Authorization"},
		ExposeHeaders:    []string{"Refresh", "Authorization"},
		AllowCredentials: true,
	}))

	// Initializing various services and repositories used in the application
	session := session.NewSession(logger, redisClient, configs.App.Auth)
	smsService := sms_service.NewSmsService(logger, &configs.SMS, TemplatesPath)

	userRepo := userRepository.NewRepository(logger, mongoDB)
	userService := service.NewUserService(logger, userRepo, session, smsService)

	chatRepo := chatRepository.NewRepository(logger, mongoDB)
	chatService := service.NewChatService(logger, chatRepo, userRepo)

	messageRepo := messageRepository.NewRepository(logger, mongoDB)
	messageService := service.NewMessageService(logger, messageRepo, chatRepo)

	// Init routes
	router.NewUserRouter(logger, app.Group("/users"), userService, chatService)

	// Init Socket Server
	socket.NewSocketService(logger, app, userService, chatService, messageService)

	// Everything almost done!
	err := app.Run(configs.App.HTTP.Address)
	if err != nil {
		panic(err)
	}
}
