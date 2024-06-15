package repository

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/kavkaco/Kavka-Core/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	db          *mongo.Database
	redisClient *redis.Client
)

func TestMain(m *testing.M) {
	database.GetMongoDBTestInstance(func(_db *mongo.Database) {
		database.GetRedisTestInstance(func(_redisClient *redis.Client) {
			fmt.Print("\n")

			db = _db
			redisClient = _redisClient

			m.Run()
		})
	})
}
