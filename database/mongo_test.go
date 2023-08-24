package database

import (
	"Kavka/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMongo(t *testing.T) {
	// Load configs
	configs := config.Read()

	// Establish connection
	mongoClient, connErr := GetMongoDBInstance(configs.Mongo)

	assert.NoError(t, connErr)
	assert.NotEmpty(t, mongoClient)
}
