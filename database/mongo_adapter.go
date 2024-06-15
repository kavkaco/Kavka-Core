package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ory/dockertest/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoLock     = &sync.Mutex{}
	mongoInstance *mongo.Database
)

var (
	UsersCollection    = "users"
	ChatsCollection    = "chats"
	MessagesCollection = "messages"
	AuthCollection     = "user_auth"
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

		ConfigureCollections(mongoInstance)
	}

	return mongoInstance, nil
}

func ConfigureCollections(db *mongo.Database) {
	db.Collection(UsersCollection).Indexes().CreateOne(context.Background(), mongo.IndexModel{ //nolint
		Keys:    bson.D{{Key: "email", Value: 1}},
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

func GetMongoDBTestInstance(callback func(db *mongo.Database)) {
	var client *mongo.Client
	var db *mongo.Database

	dockerContainerEnvVariables := []string{
		"MONGO_INITDB_ROOT_USERNAME=test",
		"MONGO_INITDB_ROOT_PASSWORD=test",
	}

	err := os.Setenv("ENV", "test")
	if err != nil {
		log.Fatalf("Could not set the environment variable to test: %s", err)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.Run("mongo", "latest", dockerContainerEnvVariables)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Kill the container
	defer func() {
		if err = pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	err = pool.Retry(func() error {
		client, err = mongo.Connect(context.TODO(),
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://test:test@localhost:%s", resource.GetPort("27017/tcp")),
			),
		)
		if err != nil {
			return err
		}

		return client.Ping(context.TODO(), nil)
	})
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %s", err)
	}

	db = client.Database("test")

	ConfigureCollections(db)

	ipAddr := resource.Container.NetworkSettings.IPAddress
	fmt.Printf("Docker container network ip address: %s\n\n", ipAddr)

	callback(db)

	if err = client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}
