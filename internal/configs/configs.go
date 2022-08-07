package configs

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type App struct {
	ListenPort int `yaml:"LISTEN_PORT"`
}

func ParseAppConfig(path string) (*App, error) {
	var configs App
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(data, &configs)

	return &configs, nil
}

type Redis struct {
	Host     string `yaml:"HOST"`
	Username string `yaml:"USERNAME"`
	Password string `yaml:"PASSWORD"`
	Port     int    `yaml:"PORT"`
	Database int    `yaml:"DB_NAME"`
}

func ParseRedisConfigs(path string) (*Redis, error) {
	var configs Redis
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(data, &configs)

	return &configs, nil
}

type Mongo struct {
	Host         string `yaml:"HOST"`
	Username     string `yaml:"USERNAME"`
	Password     string `yaml:"PASSWORD"`
	Port         int    `yaml:"PORT"`
	DatabaseName string `yaml:"DB_NAME"`
}

func ParseMongoConfigs(path string) (*Mongo, error) {
	var configs Mongo
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(data, &configs)

	return &configs, nil
}

type Smtp struct {
	Host     string `yaml:"HOST"`
	Port     int    `yaml:"PORT"`
	Email    string `yaml:"EMAIL"`
	Password string `yaml:"PASSWORD"`
}

func ParseSmtpConfigs(path string) (*Smtp, error) {
	var configs Smtp
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(data, &configs)

	return &configs, nil
}

func GetAllowedOrigins(path string) string {
	data, err := os.ReadFile(path)

	if err != nil {
		fmt.Println(err)
	}

	return string(data)
}
