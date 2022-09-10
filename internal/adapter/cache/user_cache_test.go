package cache

import (
	"Kavka/config"
	"Kavka/database"

	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const CONFIG_PATH = "/../../../config/configs.yml"
const RF_TOKEN = "sample_token"
const SAMPLE_USERNAME = "sample_username"

// NOTE - need more review & refactor probably
func TestSetAndGetToken(t *testing.T) {
	var wd, wdErr = os.Getwd()
	assert.Empty(t, wdErr)

	var configs, configsErr = config.Read(wd + CONFIG_PATH)
	assert.Empty(t, configsErr)

	var redisClient = database.GetRedisDBInstance(configs.Redis)

	var userCacheRepository = NewUserCacheRepository(redisClient, configs.App.Auth)

	// Test Set
	setErr := userCacheRepository.SetRefreshToken(RF_TOKEN, SAMPLE_USERNAME)
	assert.Empty(t, setErr)

	// Test Get
	tokenData, getErr := userCacheRepository.GetRefreshToken(RF_TOKEN)
	assert.Empty(t, getErr)
	assert.Equal(t, tokenData, &UserCache_TokenData{Username: SAMPLE_USERNAME})
}

// FIXME - real rfToken not used in this test
// var jwtManager = auth.NewJwtManager(configs.App.Auth)
// rfToken, rfTokenErr := jwtManager.GenerateRefreshToken()
// assert.Empty(t, rfTokenErr)
