package main

import (
	"Tahagram/configs"
	"Tahagram/lib"
	"Tahagram/routers"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

var AppConfigs configs.AppConfigs

func main() {
	wd, _ := os.Getwd()

	appConfigs, appConfigsErr := configs.ParseAppConfig(wd + "/configs/configs.yml")
	if appConfigsErr != nil {
		log.Fatal("Error in parsing app configs")
	}

	app := fiber.New(
		fiber.Config{
			ErrorHandler: lib.ErrorHandler,
			AppName:      "Tahagram",
			ServerHeader: "Express",
			Prefork:      true,
		},
	)

	routers.InitUsers(app)

	log.Fatal(
		app.Listen(
			fmt.Sprintf(
				"%s:%d",
				"0.0.0.0",
				appConfigs.ListenPort,
			),
		),
	)
}
