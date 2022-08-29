package database

import (
	"Kavka/configs"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBConnectionString(host string, port int, username string, password string) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)
}

func InitMongoDB(mongoConfigs configs.Mongo) (*mongo.Database, error) {
	connectionString := NewMongoDBConnectionString(
		mongoConfigs.Host,
		mongoConfigs.Port,
		mongoConfigs.Username,
		mongoConfigs.Password,
	)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	database := client.Database(mongoConfigs.DBName)

	return database, nil
}
