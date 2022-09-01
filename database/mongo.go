package database

import (
	"Kavka/config"
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoLock     = &sync.Mutex{}
	mongoInstance *mongo.Database
)

func NewMongoDBConnectionString(host string, port int, username string, password string) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)
}

func GetMongoDBInstance(mongoConfigs config.Mongo) (*mongo.Database, error) {
	if mongoInstance == nil {
		mongoLock.Lock()
		defer mongoLock.Unlock()

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

		mongoInstance = client.Database(mongoConfigs.DBName)
	}

	return mongoInstance, nil
}
