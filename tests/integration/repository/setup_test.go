package repository

import (
	"testing"

	"github.com/kavkaco/Kavka-Core/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

func TestMain(m *testing.M) {
	database.GetMongoDBTestInstance(func(_db *mongo.Database) {
		db = _db

		m.Run()
	})
}
