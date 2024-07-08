package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

var ProjectRootPath = ConfigsDirPath() + "/../"

type Env int

const (
	Development Env = iota
	Production
	Test
)

var (
	CurrentEnv Env = Development
	filename   string
)

type (
	Config struct {
		AppName string `koanf:"app_name"`
		Mongo   Mongo  `koanf:"mongo"`
		Redis   Redis  `koanf:"redis"`
		Email   Email  `koanf:"email"`
		MinIO   MinIO  `koanf:"minio"`
		HTTP    HTTP   `koanf:"http"`
		Auth    Auth   `koanf:"auth"`
	}

	Auth struct {
		SecretKey string `koanf:"secret"`
	}

	HTTP struct {
		Host string `koanf:"host"`
		Port int    `koanf:"port"`
		Cors Cors   `koanf:"cors"`
	}

	Cors struct {
		AllowOrigins []string `koanf:"allow_origins"`
	}

	Redis struct {
		Host     string `koanf:"host"`
		Username string `koanf:"username"`
		Password string `koanf:"password"`
		Port     int    `koanf:"port"`
		DB       int    `koanf:"db"`
	}

	Mongo struct {
		Host     string `koanf:"host"`
		Username string `koanf:"username"`
		Password string `koanf:"password"`
		Port     int    `koanf:"port"`
		DBName   string `koanf:"db_name"`
	}

	MinIO struct {
		Endpoint  string `koanf:"endpoint"`
		AccessKey string `koanf:"access_key"`
		SecretKey string `koanf:"secret_key"`
	}
	Email struct {
		SenderEmail string `koanf:"sender_email"`
		Password    string `koanf:"password"`
		Host        string `koanf:"host"`
		Port        string `koanf:"port"`
	}
)

func ConfigsDirPath() string {
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		panic("Error in generating env dir")
	}

	return filepath.Dir(f)
}

func Read() *Config {
	// Load KAVKA ENV
	env := strings.ToLower(os.Getenv("KAVKA_ENV"))

	if len(strings.TrimSpace(env)) == 0 || env == "development" {
		CurrentEnv = Development
		filename = "config.development.yml"
	} else if env == "production" {
		CurrentEnv = Production
		filename = "config.production.yml"
	} else if env == "test" {
		CurrentEnv = Test
		filename = "config.test.yml"
	} else {
		panic(errors.New("Invalid env value set for variable KAVKA_ENV: " + env))
	}

	// Load YAML configs
	k := koanf.New(ConfigsDirPath())
	if err := k.Load(file.Provider(fmt.Sprintf("%s/%s", ConfigsDirPath(), filename)), yaml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	config := &Config{}
	if err := k.Unmarshal("", config); err != nil {
		log.Fatalf("error unmarshaling config: %v", err)
	}

	// Load Jwt Secret Keys
	secretData, secretErr := os.ReadFile(ConfigsDirPath() + "/jwt_secret.pem")
	if secretErr != nil {
		panic(secretErr)
	}

	config.Auth.SecretKey = strings.TrimSpace(string(secretData))

	return config
}
