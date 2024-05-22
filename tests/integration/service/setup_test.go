package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/kavkaco/Kavka-Core/database"
	"github.com/ory/dockertest/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoEnvVariables = []string{
	"MONGO_INITDB_ROOT_USERNAME=test",
	"MONGO_INITDB_ROOT_PASSWORD=test",
}

var redisEnvVariables = []string{}

var (
	db          *mongo.Database
	redisClient *redis.Client
)

func getMongoTestInstance(pool *dockertest.Pool) (*mongo.Client, *dockertest.Resource) {
	var c *mongo.Client

	resource, err := pool.Run("mongo", "latest", mongoEnvVariables)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	err = pool.Retry(func() error {
		c, err = mongo.Connect(context.TODO(),
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://test:test@localhost:%s", resource.GetPort("27017/tcp")),
			),
		)
		if err != nil {
			return err
		}

		return c.Ping(context.TODO(), nil)
	})
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %s", err)
	}

	return c, resource
}

func getRedisTestInstance(pool *dockertest.Pool) (*redis.Client, *dockertest.Resource) {
	var c *redis.Client

	resource, err := pool.Run("redis", "alpine", redisEnvVariables)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	err = pool.Retry(func() error {
		ipAddr := resource.Container.NetworkSettings.IPAddress + ":6379"

		c = redis.NewClient(&redis.Options{
			Addr: ipAddr,
			DB:   0,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Could not connect to Redis: %s", err)
	}

	return c, resource
}

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

	mc, mongoResource := getMongoTestInstance(pool)
	rc, redisResource := getRedisTestInstance(pool)

	db = mc.Database("test")
	database.ConfigureCollections(db)

	redisClient = rc

	code := m.Run()

	// Kill the containers
	if err = pool.Purge(mongoResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	if err = pool.Purge(redisResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
