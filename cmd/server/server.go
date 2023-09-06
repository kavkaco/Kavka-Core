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

	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Response().Header.Add("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Refresh, Authorization")
		return ctx.Next()
	})

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
