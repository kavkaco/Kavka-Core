package repository

import (
	"os"
	"testing"

	"github.com/kavkaco/Kavka-Core/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	dbClient *mongo.Client
	db       *mongo.Database
)

func TestMain(m *testing.M) {
	database.GetMongoDBTestInstance(func(db *mongo.Database) {
		code := m.Run()

		os.Exit(code)
	})
}
