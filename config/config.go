package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

var (
	EnvItems   = []string{"devel", "prod"}
	CurrentEnv string
)

type (
	IConfig struct {
		App   App   `yaml:"app"`
		Mongo Mongo `yaml:"mongo"`
		Redis Redis `yaml:"redis"`
		SMS   `yaml:"sms"`
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
		DB       int    `yaml:"db"`
		DBName   string `yaml:"db_name"`
	}
	MinIOCredentials struct {
		Endpoint  string `yaml:"endpoint"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
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
		env = EnvItems[0]
	} else if slices.Contains(EnvItems, env) {
		CurrentEnv = env
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

	return cfg
}
