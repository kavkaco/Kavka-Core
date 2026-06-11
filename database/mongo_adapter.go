package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

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

		clientOpts := options.Client().
			ApplyURI(uri).
			SetMaxPoolSize(500).
			SetMinPoolSize(50).
			SetMaxConnIdleTime(30 * time.Second).
			SetMaxConnecting(50).
			SetSocketTimeout(10 * time.Second).
			SetConnectTimeout(5 * time.Second).
			SetHeartbeatInterval(10 * time.Second)

		client, err := mongo.Connect(context.Background(), clientOpts)
		if err != nil {
			return nil, err
		}

		// Send a ping to confirm a successful connection
		var result bson.M
		if err := client.Database("test").RunCommand(context.Background(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
			return nil, err
		}

		mongoInstance = client.Database(dbName)

		ConfigureCollections(mongoInstance)
	}

	return mongoInstance, nil
}

func ConfigureCollections(db *mongo.Database) {
	handleError := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	ctx := context.Background()

	// Users indexes

	_, err := db.Collection(UsersCollection).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"user_id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"username": 1},
			Options: options.Index().SetUnique(true),
		},
	})
	handleError(err)

	_, err = db.Collection(UsersCollection).Indexes().CreateMany(ctx, []mongo.IndexModel{ //nolint
		{
			Keys: bson.D{
				{Key: "name", Value: "text"},
				{Key: "email", Value: "text"},
				{Key: "username", Value: "text"},
				{Key: "last_name", Value: "text"},
			},
		},
	})
	handleError(err)

	// Chats indexes

	_, err = db.Collection(ChatsCollection).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "chat_detail.username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "chat_detail.title", Value: "text"},
				{Key: "chat_detail.username", Value: "text"},
			},
			Options: options.Index(),
		},
	})
	handleError(err)

	// Messages indexes

	_, err = db.Collection(MessagesCollection).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "chat_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "chat_id", Value: 1},
				{Key: "message_id", Value: 1},
			},
		},
	})
	handleError(err)

	// User auth indexes

	_, err = db.Collection(AuthCollection).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	})
	handleError(err)

	// Messages V2 indexes

	_, err = db.Collection("messages_v2").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "chat_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "chat_id", Value: 1},
				{Key: "_id", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "chat_id", Value: 1}},
		},
	})
	handleError(err)
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
		client, err = mongo.Connect(context.Background(),
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://test:test@localhost:%s", resource.GetPort("27017/tcp")),
			),
		)
		if err != nil {
			return err
		}

		return client.Ping(context.Background(), nil)
	})
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %s", err)
	}

	db = client.Database("test")

	ConfigureCollections(db)

	ipAddr := resource.Container.NetworkSettings.IPAddress + ":27017"

	fmt.Printf("Docker mongodb container network ip address: %s\n", ipAddr)

	callback(db)

	if err = client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}
