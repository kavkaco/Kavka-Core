package main

import (
	"Kavka/app/router"
	"Kavka/config"
	"Kavka/database"
	repository "Kavka/internal/repository/user"
	"Kavka/internal/service"
	"Kavka/pkg/session"
	"Kavka/utils/sms_otp"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	// Define paths
	var (
		TEMPLATES_PATH = config.ProjectRootPath + "/app/views/mail/"
	)

	// Load Configs
	configs := config.Read()

	// Init MongoDB
	mongoDB, mongoErr := database.GetMongoDBInstance(configs.Mongo)
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
	smsOtp := sms_otp.NewSMSOtpService(&configs.SMS, TEMPLATES_PATH)

	userRepo := repository.NewUserRepository(mongoDB)
	userService := service.NewUserService(userRepo, session, smsOtp)
	router.NewUserRouter(app.Group("/users"), userService)

	// Init Socket Server
	// socket.NewSocketService(app, userService)

	// Everything almost done!
	app.Run(configs.App.HTTP.Address)
}
