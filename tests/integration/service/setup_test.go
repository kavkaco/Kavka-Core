package service

import (
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/infra/stream"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	db          *mongo.Database
	redisClient *redis.Client
	natsConn    *nats.Conn
)

func TestMain(m *testing.M) {
	database.GetMongoDBTestInstance(func(_db *mongo.Database) {
		database.GetRedisTestInstance(func(_redisClient *redis.Client) {
			stream.GetNATSTestInstance(func(_natsConn *nats.Conn) {
				db = _db
				redisClient = _redisClient
				natsConn = _natsConn

				m.Run()
			})
		})
	})
}
