package main

import (
	"Tahagram/api"
	"Tahagram/configs"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	mainPath, _ := os.Getwd()
	configs, configsErr := configs.ParseConfig(mainPath + "/configs/configs.yml")
	if configsErr != nil {
		log.Fatal(configsErr)
	}

	app := fiber.New(
		fiber.Config{
			ErrorHandler: api.ErrorHandler,
			AppName:      "Tahagram",
			ServerHeader: "Express",
			Prefork:      true,
		},
	)

	log.Fatal(
		app.Listen(
			fmt.Sprintf(
				"%s:%d",
				"0.0.0.0",
				configs.ListenPort,
			),
		),
	)
}
