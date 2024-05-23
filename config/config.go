package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

type Env int

const (
	Development Env = iota
	Production
	Test
)

var CurrentEnv Env = Development

type (
	IConfig struct {
		App   App              `yaml:"app"`
		Mongo Mongo            `yaml:"mongo"`
		Redis Redis            `yaml:"redis"`
		Email Email            `yaml:"email"`
		MinIO MinIOCredentials `yaml:"minio"`
	}
	App struct {
		Name   string `yaml:"name"`
		HTTP   HTTP   `yaml:"http"`
		Server Server `yaml:"server"`
		Auth   Auth   `yaml:"auth"`
	}
	HTTP struct {
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		Address string `yaml:"address"`
	}
	Auth struct {
		SECRET string `yaml:"secret"`
	}
	Server struct {
		CORS CORS `yaml:"cors"`
	}
	CORS struct {
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
	MinIOCredentials struct {
		Endpoint  string `yaml:"endpoint"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
	}
	Email struct{}
)

var ProjectRootPath = ConfigsDirPath() + "/../"

func ConfigsDirPath() string {
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		panic("Error in generating env dir")
	}

	return filepath.Dir(f)
}

func Read() *IConfig {
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
	var cfg *IConfig

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

	cfg.App.Auth.SECRET = strings.TrimSpace(string(secretData))

	return cfg
}
