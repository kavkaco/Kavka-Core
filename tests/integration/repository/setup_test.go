package repository

import (
	"fmt"
	"testing"

	"github.com/kavkaco/Kavka-Core/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	db *mongo.Database
)

func TestMain(m *testing.M) {
	database.GetMongoDBTestInstance(func(_db *mongo.Database) {
		fmt.Print("\n")

		db = _db

		m.Run()
	})
}
