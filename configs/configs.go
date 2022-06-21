package configs

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfigs struct {
	ListenPort int `yaml:"LISTEN_PORT"`
}

func ParseAppConfig(path string) (*AppConfigs, error) {
	var configs AppConfigs
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(data, &configs)

	return &configs, nil
}

type RedisConfigs struct {
	Host     string `yaml:"HOST"`
	Username string `yaml:"USERNAME"`
	Password string `yaml:"PASSWORD"`
	Port     int    `yaml:"PORT"`
	Database int    `yaml:"DB_NAME"`
}

func ParseRedisConfigs(path string) (*RedisConfigs, error) {
	var configs RedisConfigs
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(data, &configs)

	return &configs, nil
}

type MongoConfigs struct {
	Host         string `yaml:"HOST"`
	Username     string `yaml:"USERNAME"`
	Password     string `yaml:"PASSWORD"`
	Port         int    `yaml:"PORT"`
	DatabaseName string `yaml:"DB_NAME"`
}

func ParseMongoConfigs(path string) (*MongoConfigs, error) {
	var configs MongoConfigs
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(data, &configs)

	return &configs, nil
}

func GetAllowOrigins() string {
	wd, _ := os.Getwd()

	data, err := os.ReadFile(wd + "/configs/allow_origins")

	if err != nil {
		fmt.Println(err)
	}

	return string(data)
}
