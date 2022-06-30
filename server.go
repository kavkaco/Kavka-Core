package main

import (
	"Tahagram/configs"
	"Tahagram/database"
	"Tahagram/logs"
	"Tahagram/routers"
	"Tahagram/websocket"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var AppConfigs configs.AppConfigs
var MongoConfigs configs.MongoConfigs

func ParseConfigs() {
	wd, _ := os.Getwd()

	appConfigs, appConfigsErr := configs.ParseAppConfig(wd + "/configs/configs.yml")
	if appConfigsErr != nil {
		fmt.Println("Error in parsing app configs")
		os.Exit(1)
	}

	mongoConfigs, mongoConfigsErr := configs.ParseMongoConfigs(wd + "/configs/mongo.yml")
	if mongoConfigsErr != nil {
		fmt.Println("Error in parsing mongodb configs")
	}

	AppConfigs = *appConfigs
	MongoConfigs = *mongoConfigs
}

func main() {
	ParseConfigs()

	database.EstablishConnection(MongoConfigs)

	app := fiber.New(
		fiber.Config{
			AppName:      "Tahagram",
			ServerHeader: "Fiber",
			Prefork:      true,
		},
	)

	app.Use(cors.New(
		cors.Config{
			AllowOrigins:     configs.GetAllowOrigins(),
			AllowCredentials: true,
		},
	))

	logs.InitLogger(app)

	api := app.Group("/api")
	routers.InitUsers(api)

	websocket.InitWebSocket(app)

	log.Fatal(
		app.Listen(
			fmt.Sprintf(
				"%s:%d",
				"0.0.0.0",
				AppConfigs.ListenPort,
			),
		),
	)
}