package database

import (
	"Tahagram/configs"
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database
var ctx = context.TODO()

var (
	UsersCollection *mongo.Collection
)

func MakeConnectionString(host string, port int, username string, password string) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)
}

func EstablishConnection() {
	wd, _ := os.Getwd()

	mongoConfigs, mongoConfigsErr := configs.ParseMongoConfigs(wd + "/configs/mongo.yml")
	if mongoConfigsErr != nil {
		fmt.Println("Error in parsing mongodb configs")
	}

	connectionString := MakeConnectionString(
		mongoConfigs.Host,
		mongoConfigs.Port,
		mongoConfigs.Username,
		mongoConfigs.Password,
	)

	client, clientErr := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if clientErr != nil {
		fmt.Println("Error in connecting to mongo database:")
		fmt.Println(clientErr)
		os.Exit(1)
	} else {
		fmt.Println("Successfully connected to mongo database")
	}

	MongoDB = client.Database(mongoConfigs.DatabaseName)
	UsersCollection = MongoDB.Collection("users")
}
