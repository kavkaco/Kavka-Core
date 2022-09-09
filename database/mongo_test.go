package database

import (
	"Kavka/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMongo(t *testing.T) {
	// Read configs

	cfg, err := config.Read("../config/configs.yml")

	assert.NoError(t, err)

	// Establish connection

	mongoClient, connErr := GetMongoDBInstance(cfg.Mongo)

	assert.NoError(t, connErr)
	assert.NotEmpty(t, mongoClient)
}
