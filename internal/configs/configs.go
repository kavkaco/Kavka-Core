package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

const CONFIG_PATH = "./app/configs/configs.yml"

type App struct {
	Name  string `yaml:"NAME"`
	HTTP  HTTP   `yaml:"HTTP"`
	Fiber Fiber  `yaml:"FIBER"`
}
type HTTP struct {
	Host    string `yaml:"HOST"`
	Port    int    `yaml:"PORT"`
	Address string `yaml:"ADDRESS"`
}

type Fiber struct {
	ServerHeader string `yaml:"SERVER_HEADER"`
	Prefork      bool   `yaml:"PREFORK"`
	CORS         CORS   `yaml:"CORS"`
}
type CORS struct {
	AllowOrigins     string `yaml:"ALLOW_ORIGINS"`
	AllowCredentials bool   `yaml:"ALLOW_CREDENTIALS"`
}
type Redis struct {
	Host     string `yaml:"HOST"`
	Username string `yaml:"USERNAME"`
	Password string `yaml:"PASSWORD"`
	Port     int    `yaml:"PORT"`
}
type Mongo struct {
	Host     string `yaml:"HOST"`
	Username string `yaml:"USERNAME"`
	Password string `yaml:"PASSWORD"`
	Port     int    `yaml:"PORT"`
	DBName   string `yaml:"DB_NAME"`
}

type SMTP struct {
	Host     string `yaml:"HOST"`
	Port     int    `yaml:"PORT"`
	Email    string `yaml:"EMAIL"`
	Password string `yaml:"PASSWORD"`
}
type Config struct {
	App   App   `yaml:"APP"`
	Mongo Mongo `yaml:"MONGO"`
	Redis Redis `yaml:"REDIS"`
	SMTP  SMTP  `yaml:"SMTP"`
}

func Parse() (*Config, error) {
	var configs Config
	data, err := os.ReadFile(CONFIG_PATH)
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(data, &configs)

	return &configs, nil
}
