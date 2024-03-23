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
	"github.com/kavkaco/Kavka-Core/pkg/session"
	"github.com/kavkaco/Kavka-Core/pkg/sms_service"
)

func main() {
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

	// ----- Init Services -----
	session := session.NewSession(redisClient, configs.App.Auth)
	smsService := sms_service.NewSmsService(&configs.SMS, TemplatesPath)

	userRepo := userRepository.NewRepository(mongoDB)
	userService := service.NewUserService(userRepo, session, smsService)

	chatRepo := chatRepository.NewRepository(mongoDB)
	chatService := service.NewChatService(chatRepo, userRepo)

	messageRepo := messageRepository.NewRepository(mongoDB)
	messageRepository := service.NewMessageService(messageRepo, chatRepo)

	// Init routes
	router.NewUserRouter(app.Group("/users"), userService, chatService)

	// Init Socket Server
	socket.NewSocketService(app, userService, chatService, messageRepository)

	// Everything almost done!
	err := app.Run(configs.App.HTTP.Address)
	if err != nil {
		panic(err)
	}
}
