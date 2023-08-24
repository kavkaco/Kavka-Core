package main

import (
	"Kavka/app/middleware"
	"Kavka/app/router"
	"Kavka/app/socket"
	"Kavka/config"
	"Kavka/database"
	repository "Kavka/internal/repository/user"
	"Kavka/internal/service"
	"Kavka/pkg/session"
	"Kavka/utils/sms_otp"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// const (
// 	// YAML config file path
// 	CONFIG_PATH string = "./config/configs.yml"
// 	// Define templates path (used by: SmsOtpService)
// 	templatesPath = "./app/views/mail/"
// )

func main() {
	config.SetCwd()

	// Define paths
	var (
		CONFIG_PATH    string = config.Cwd + "/config/configs.yml"
		TEMPLATES_PATH        = config.Cwd + "/app/views/mail/"
	)

	// Load Configs
	configs, configErr := config.Read(CONFIG_PATH)
	if configErr != nil {
		panic(configErr)
	}

	// Init MongoDB
	mongoDB, mongoErr := database.GetMongoDBInstance(configs.Mongo)
	if mongoErr != nil {
		panic(mongoErr)
	}

	// Init RedisDB
	redisClient := database.GetRedisDBInstance(configs.Redis)

	// Init WebServer
	app := fiber.New(
		fiber.Config{
			AppName:      configs.App.Name,
			Prefork:      configs.App.Fiber.Prefork,
			ErrorHandler: middleware.ErrorHandler,
		},
	)

	// Config WebServer
	app.Use(cors.New(
		cors.Config{
			AllowOrigins:     configs.App.Fiber.CORS.AllowOrigins,
			AllowCredentials: configs.App.Fiber.CORS.AllowCredentials,
		},
	))

	// ----- Init Services -----
	session := session.NewSession(redisClient, configs.App.Auth)
	smsOtp := sms_otp.NewSMSOtpService(&configs.SMS, TEMPLATES_PATH)

	userRepo := repository.NewUserRepository(mongoDB)
	userService := service.NewUserService(userRepo, session, smsOtp)
	router.NewUserRouter(app.Group("/users"), userService)

	// Init Socket Server
	socket.NewSocketService(app, userService)

	// Everything almost done!
	log.Fatal(
		app.Listen(configs.App.HTTP.Address),
	)
}
