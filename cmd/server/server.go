package main

import (
	"Kavka/config"
	"Kavka/database"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const CONFIG_PATH string = "./config/configs.yml"

func main() {
	configs, configErr := config.Read(CONFIG_PATH)
	if configErr != nil {
		panic(configErr)
	}

	_, mongoErr := database.GetMongoDBInstance(configs.Mongo)
	if mongoErr != nil {
		panic(mongoErr)
	}

	// TODO - move
	app := fiber.New(
		fiber.Config{
			AppName:      configs.App.Name,
			ServerHeader: configs.App.Fiber.ServerHeader,
			Prefork:      configs.App.Fiber.Prefork,
		},
	)

	app.Use(cors.New(
		cors.Config{
			AllowOrigins:     configs.App.Fiber.CORS.AllowOrigins,
			AllowCredentials: configs.App.Fiber.CORS.AllowCredentials,
		},
	))

	log.Fatal(
		app.Listen(configs.App.HTTP.Address),
	)
}
