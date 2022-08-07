package database

import (
	"Kavka/internal/configs"
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

func InitMongoDB(mongoConfigs configs.Mongo) {
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
	}

	MongoDB = client.Database(mongoConfigs.DatabaseName)

	UsersCollection = MongoDB.Collection("userDB")
}
