package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

var (
	ENV_ITEMS = []string{"devel", "prod"}
	ENV       string
)

type (
	IConfig struct {
		App   App   `json:"app" yaml:"APP"`
		Mongo Mongo `json:"mongo" yaml:"MONGO"`
		Redis Redis `json:"redis" yaml:"REDIS"`
		SMS   `json:"sms" yaml:"SMS"`
		MinIOCredentials
	}
	App struct {
		Name   string `json:"name" yaml:"NAME"`
		HTTP   HTTP   `json:"http" yaml:"HTTP"`
		Server Server `json:"server" yaml:"SERVER"`
		Auth   Auth   `json:"auth" yaml:"AUTH"`
	}
	HTTP struct {
		Host    string `json:"host" yaml:"HOST"`
		Port    int    `json:"port" yaml:"PORT"`
		Address string `json:"address" yaml:"ADDRESS"`
	}
	Auth struct {
		SECRET string `json:"secret" yaml:"SECRET"`
	}
	Server struct {
		CORS CORS `json:"cors" yaml:"CORS"`
	}
	CORS struct {
		AllowOrigins string `json:"allow_origins" yaml:"ALLOW_ORIGINS"`
	}
	Redis struct {
		Host     string `json:"host" yaml:"HOST"`
		Username string `json:"username" yaml:"USERNAME"`
		Password string `json:"password" yaml:"PASSWORD"`
		Port     int    `json:"port" yaml:"PORT"`
		DB       int    `json:"db" yaml:"DB"`
	}
	Mongo struct {
		Host     string `json:"host" yaml:"HOST"`
		Username string `json:"username" yaml:"USERNAME"`
		Password string `json:"password" yaml:"PASSWORD"`
		Port     int    `json:"port" yaml:"PORT"`
		DBName   string `json:"db_name" yaml:"DB_NAME"`
	}
	MinIOCredentials struct {
		Endpoint  string `json:"ENDPOINT"   yaml:"ENDPOINT"`
		AccessKey string `json:"ACCESS_KEY" yaml:"ACCESS_KEY"`
		SecretKey string `json:"SECRET_KEY" yaml:"SECRET_KEY"`
	}
	// TODO - Add sms-service's configs.
	SMS struct{}
)

const _ = "/config/configs.yml"

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
	env := os.Getenv("ENV")
	if len(strings.TrimSpace(env)) == 0 {
		ENV = ENV_ITEMS[0]
	} else if slices.Contains(ENV_ITEMS, env) {
		ENV = env
	} else {
		panic(errors.New("Invalid ENV key: " + env))
	}

	// Load YAML configs
	var cfg *IConfig

	data, readErr := os.ReadFile(ConfigsDirPath() + "/configs.yml")
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

	// Load MinIO credentials
	filename := "minio-credentials.json"
	credFile, credErr := os.ReadFile(ConfigsDirPath() + "/" + filename)
	if credErr != nil {
		panic(credErr)
	}

	var cred MinIOCredentials

	jsonErr := json.Unmarshal(credFile, &cred)
	if jsonErr != nil {
		panic(jsonErr)
	}

	cfg.MinIOCredentials = cred

	return cfg
}
