package database

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoLock     = &sync.Mutex{}
	mongoInstance *mongo.Database
)

var (
	UsersCollection = "users"
	ChatsCollection = "chats"
)

func NewMongoDBConnectionString(host string, port int, username string, password string) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port) //nolint
}

func GetMongoDBInstance(uri, dbName string) (*mongo.Database, error) {
	if mongoInstance == nil {
		mongoLock.Lock()
		defer mongoLock.Unlock()

		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err != nil {
			return nil, err
		}

		// Send a ping to confirm a successful connection
		var result bson.M
		if err := client.Database("test").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
			return nil, err
		}

		mongoInstance = client.Database(dbName)

		collectionsConfigurations(mongoInstance)
	}

	return mongoInstance, nil
}

func collectionsConfigurations(db *mongo.Database) {
	db.Collection(UsersCollection).Indexes().CreateOne(context.Background(), mongo.IndexModel{ //nolint
		Keys:    bson.D{{Key: "phone", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	db.Collection(UsersCollection).Indexes().CreateOne(context.Background(), mongo.IndexModel{ //nolint
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
}

func IsDuplicateKeyError(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}

func IsRowExistsError(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 121 {
				return true
			}
		}
	}
	return false
}
