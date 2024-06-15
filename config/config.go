package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

var ProjectRootPath = ConfigsDirPath() + "/../"

type Env int

const (
	Development Env = iota
	Production
	Test
)

var CurrentEnv Env = Development

type (
	Config struct {
		AppName string `yaml:"app_name"`
		Mongo   Mongo  `yaml:"mongo"`
		Redis   Redis  `yaml:"redis"`
		Email   Email  `yaml:"email"`
		MinIO   MinIO  `yaml:"minio"`
		HTTP    HTTP   `yaml:"http"`
		Auth    Auth   `yaml:"auth"`
	}
	Auth struct {
		SecretKey string `yaml:"secret"`
	}
	HTTP struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		Cors Cors   `yaml:"cors"`
	}
	Cors struct {
		AllowOrigins string `yaml:"allow_origins"`
	}
	Redis struct {
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Port     int    `yaml:"port"`
		DB       int    `yaml:"db"`
	}
	Mongo struct {
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Port     int    `yaml:"port"`
		DBName   string `yaml:"db_name"`
	}
	MinIO struct {
		Endpoint  string `yaml:"endpoint"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
	}
	Email struct{}
)

func ConfigsDirPath() string {
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		panic("Error in generating env dir")
	}

	return filepath.Dir(f)
}

func Read() *Config {
	// Load ENV
	env := strings.ToLower(os.Getenv("ENV"))
	if len(strings.TrimSpace(env)) == 0 || env == "development" {
		CurrentEnv = Development
	} else if env == "production" {
		CurrentEnv = Production
	} else if env == "test" {
		CurrentEnv = Test
	} else {
		panic(errors.New("Invalid ENV: " + env))
	}

	// Load YAML configs
	var cfg *Config

	data, readErr := os.ReadFile(ConfigsDirPath() + "/config.yml")
	if readErr != nil {
		panic(readErr)
	}

	parseErr := yaml.Unmarshal(data, &cfg)
	if parseErr != nil {
		panic(parseErr)
	}

	// Load JwtSecret keys
	secretData, secretErr := os.ReadFile(ConfigsDirPath() + "/jwt_secret.pem")
	if secretErr != nil {
		panic(secretErr)
	}

	cfg.Auth.SecretKey = strings.TrimSpace(string(secretData))

	return cfg
}
