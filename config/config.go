package config

import (
	"errors"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

var Cwd string
var ENV_ITEMS = []string{"devel", "prod"}
var ENV string

type (
	IConfig struct {
		App   App   `yaml:"APP"`
		Mongo Mongo `yaml:"MONGO"`
		Redis Redis `yaml:"REDIS"`
		SMS   `yaml:"SMS"`
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
		SECRET             string
		OTP_EXPIRE_SECONDS time.Duration `json:"OTP_EXPIRE_SECONDS"`
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
	// TODO - Add sms-service's configs
	SMS struct{}
)

func SetCwd() {
	wd, _ := os.Getwd()
	Cwd = wd
}

func Read(fileName string) (IConfig, error) {
	// Load ENV
	env := os.Getenv("ENV")
	if len(strings.TrimSpace(env)) == 0 {
		ENV = ENV_ITEMS[0]
	} else if slices.Contains(ENV_ITEMS, env) {
		ENV = env
	} else {
		return IConfig{}, errors.New("Invalid ENV key: " + env)
	}

	// Load YAML configs
	var cfg *IConfig

	data, readErr := os.ReadFile(fileName)
	if readErr != nil {

	}

	parseErr := yaml.Unmarshal(data, &cfg)
	if parseErr != nil {
		return IConfig{}, parseErr
	}

	// Load RSA keys
	secretData, secretErr := os.ReadFile(Cwd + "/config/jwt_secret.pem")
	if secretErr != nil {
		return IConfig{}, secretErr
	}
	cfg.App.Auth.SECRET = strings.TrimSpace(string(secretData))

	return *cfg, nil
}
