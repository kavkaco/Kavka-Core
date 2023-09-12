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
		App   App   `yaml:"app"`
		Mongo Mongo `yaml:"mongo"`
		Redis Redis `yaml:"redis"`
		SMS   `yaml:"sms"`
		MinIOCredentials
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
		DB       int    `yaml:"db"`
		DBName   string `yaml:"db_name"`
	}
	MinIOCredentials struct {
		Endpoint  string `json:"endpoint"   yaml:"endpoint"`
		AccessKey string `json:"access_key" yaml:"access_key"`
		SecretKey string `json:"secret_key" yaml:"secret_key"`
	}
	// TODO - Add sms-service's configs.
	SMS struct{}
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
