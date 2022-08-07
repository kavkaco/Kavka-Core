package main

import (
	"Kavka/app/database"
	"Kavka/app/routers"
	"Kavka/app/session"
	"Kavka/app/websocket"
	"Kavka/internal/configs"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const APP_CONFIG_PATH string = "./app/configs/configs.yml"

func main() {

	cfg, err := configs.Parse(APP_CONFIG_PATH)
	if err != nil {
		log.Fatal("cannot parse configs. ", err.Error())
	}

	redisClient := database.InitRedisDB(cfg.Redis)
	database.InitMongoDB(cfg.Mongo)

	app := fiber.New(
		fiber.Config{
			AppName:      cfg.App.Name,
			ServerHeader: cfg.App.Fiber.ServerHeader,
			Prefork:      cfg.App.Fiber.Prefork,
		},
	)

	app.Use(cors.New(
		cors.Config{
			AllowOrigins:     cfg.App.Fiber.CORS.AllowOrigins,
			AllowCredentials: cfg.App.Fiber.CORS.AllowCredentials,
		},
	))

	api := app.Group("/api")
	routers.InitUsers(api)
	websocket.InitWebSocket(app)
	session.InitSession(redisClient)

	log.Fatal(
		app.Listen(cfg.App.HTTP.Address),
	)
}
