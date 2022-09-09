package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type (
	IConfig struct {
		App   App   `yaml:"APP"`
		Mongo Mongo `yaml:"MONGO"`
		Redis Redis `yaml:"REDIS"`
		SMTP  SMTP  `yaml:"SMTP"`
	}
	App struct {
		Name  string `yaml:"NAME"`
		HTTP  HTTP   `yaml:"HTTP"`
		Fiber Fiber  `yaml:"FIBER"`
		Auth  Auth   `yaml:"AUTH"`
	}
	HTTP struct {
		Host    string `yaml:"HOST"`
		Port    int    `yaml:"PORT"`
		Address string `yaml:"ADDRESS"`
	}
	Auth struct {
		JWTSecretKey  string        `json:"JWT_SECRET_KEY"`
		AT_TTL_MINUTE time.Duration `json:"AT_TTL_MINUTE"`
		RF_TTL_MINUTE time.Duration `json:"RF_TTL_MINUTE"` // Refresh Token TTL minute
	}
	Fiber struct {
		ServerHeader string `yaml:"SERVER_HEADER"`
		Prefork      bool   `yaml:"PREFORK"`
		CORS         CORS   `yaml:"CORS"`
	}
	CORS struct {
		AllowOrigins     string `yaml:"ALLOW_ORIGINS"`
		AllowCredentials bool   `yaml:"ALLOW_CREDENTIALS"`
	}
	Redis struct {
		Host     string `yaml:"HOST"`
		Username string `yaml:"USERNAME"`
		Password string `yaml:"PASSWORD"`
		Port     int    `yaml:"PORT"`
		DB       int    `yaml:"DB"`
	}
	Mongo struct {
		Host     string `yaml:"HOST"`
		Username string `yaml:"USERNAME"`
		Password string `yaml:"PASSWORD"`
		Port     int    `yaml:"PORT"`
		DBName   string `yaml:"DB_NAME"`
	}

	SMTP struct {
		Host     string `yaml:"HOST"`
		Port     int    `yaml:"PORT"`
		Email    string `yaml:"EMAIL"`
		Password string `yaml:"PASSWORD"`
	}
)

func Read(fileName string) (IConfig, error) {
	var cfg *IConfig

	data, readErr := os.ReadFile(fileName)
	if readErr != nil {
		return IConfig{}, readErr
	}

	parseErr := yaml.Unmarshal(data, &cfg)
	if parseErr != nil {
		return IConfig{}, parseErr
	}

	return *cfg, nil
}
