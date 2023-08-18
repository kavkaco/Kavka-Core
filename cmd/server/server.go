package main

import (
	"Kavka/configs"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const CONFIG_PATH string = "./app/configs/configs.yml"

func main() {
	cfg, err := configs.Read(CONFIG_PATH)
	if err != nil {
		log.Fatal("cannot parse configs: ", err.Error())
	}

	// FIXME
	// mongoDB, mongoErr := database.InitMongoDB(cfg.Mongo)
	// if mongoErr != nil {
	// 	log.Fatal("MongoDB Connection Error :", mongoErr.Error())
	// }

	// TODO - move
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

	log.Fatal(
		app.Listen(cfg.App.HTTP.Address),
	)
}
