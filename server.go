package main

import (
	"Kavka/app/database"
	"Kavka/app/routers"
	"Kavka/app/session"
	"Kavka/app/websocket"
	"Kavka/internal/configs"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var AppConfigs configs.AppConfigs
var MongoConfigs configs.MongoConfigs
var RedisConfigs configs.RedisConfigs
var AllowedOrigins string

func ParseConfigs() {
	wd, _ := os.Getwd()

	appConfigs, appConfigsErr := configs.ParseAppConfig(wd + "/app/configs/configs.yml")
	if appConfigsErr != nil {
		fmt.Println("Error in parsing app configs")
		os.Exit(1)
	}

	mongoConfigs, mongoConfigsErr := configs.ParseMongoConfigs(wd + "/app/configs/mongo.yml")
	if mongoConfigsErr != nil {
		fmt.Println("Error in parsing mongodb configs")
	}

	redisConfigs, redisConfigsErr := configs.ParseRedisConfigs(wd + "/app/configs/redis.yml")
	if redisConfigsErr != nil {
		log.Fatal("Error in connecting to redis database")
	}

	allowedOrigins := configs.GetAllowedOrigins(wd + "/app/configs/allowed_origins")

	AppConfigs = *appConfigs
	MongoConfigs = *mongoConfigs
	RedisConfigs = *redisConfigs
	AllowedOrigins = allowedOrigins
}

func main() {
	ParseConfigs() // FIXME - will change

	redisClient := database.InitRedisDB(RedisConfigs)
	database.InitMongoDB(MongoConfigs)

	app := fiber.New(
		fiber.Config{
			AppName:      "Kavka",
			ServerHeader: "Fiber",
			// Prefork:      true,
		},
	)

	app.Use(cors.New(
		cors.Config{
			AllowOrigins:     AllowedOrigins,
			AllowCredentials: true,
		},
	))

	api := app.Group("/api")
	routers.InitUsers(api)
	websocket.InitWebSocket(app)
	session.InitSession(redisClient)

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
