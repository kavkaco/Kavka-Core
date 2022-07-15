package main

import (
	"Nexus/internal/configs"
	"os"
	"testing"
)

func TestAppConfigs(t *testing.T) {
	wd, _ := os.Getwd()

	appConfigs, _ := configs.ParseAppConfig(wd + "/../configs/configs.yml")

	t.Log(appConfigs)
}

func TestRedisConfigs(t *testing.T) {
	wd, _ := os.Getwd()

	redisConfigs, _ := configs.ParseRedisConfigs(wd + "/../configs/redis.yml")

	t.Log(redisConfigs)
}

func TestMongoConfigs(t *testing.T) {
	wd, _ := os.Getwd()

	mongoConfigs, _ := configs.ParseMongoConfigs(wd + "/../configs/mongo.yml")

	t.Log(mongoConfigs)
}
