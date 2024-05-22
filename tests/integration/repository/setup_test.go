package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/ory/dockertest/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DockerContainerEnvVariables = []string{
	"MONGO_INITDB_ROOT_USERNAME=test",
	"MONGO_INITDB_ROOT_PASSWORD=test",
}

var (
	dbClient *mongo.Client
	db       *mongo.Database
)

func TestMain(m *testing.M) {
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

	resource, err := pool.Run("mongo", "latest", DockerContainerEnvVariables)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	err = pool.Retry(func() error {
		dbClient, err = mongo.Connect(context.TODO(),
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://test:test@localhost:%s", resource.GetPort("27017/tcp")),
			),
		)
		if err != nil {
			return err
		}

		return dbClient.Ping(context.TODO(), nil)
	})
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %s", err)
	}

	db = dbClient.Database("test")

	database.ConfigureCollections(db)

	ipAddr := resource.Container.NetworkSettings.IPAddress
	fmt.Printf("Docker container network ip address: %s\n\n", ipAddr)

	code := m.Run()

	if err = dbClient.Disconnect(context.Background()); err != nil {
		panic(err)
	}

	// Kill the container
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
