package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

type IConfigs struct {
	App struct {
		Name string `yaml:"NAME"`

		HTTP struct {
			Host    string `yaml:"HOST"`
			Port    int    `yaml:"PORT"`
			Address string `yaml:"ADDRESS"`
		} `yaml:"HTTP"`

		Fiber struct {
			ServerHeader string `yaml:"SERVER_HEADER"`
			Prefork      bool   `yaml:"PREFORK"`
			CORS         struct {
				AllowOrigins     string `yaml:"ALLOW_ORIGINS"`
				AllowCredentials bool   `yaml:"ALLOW_CREDENTIALS"`
			} `yaml:"CORS"`
		} `yaml:"FIBER"`
	} `yaml:"APP"`

	Redis struct {
		Host     string `yaml:"HOST"`
		Username string `yaml:"USERNAME"`
		Password string `yaml:"PASSWORD"`
		Port     int    `yaml:"PORT"`
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
}

func Read(fileName string) (IConfigs, error) {
	var cfg *IConfigs

	data, readErr := os.ReadFile(fileName)
	if readErr != nil {
		return IConfigs{}, readErr
	}

	parseErr := yaml.Unmarshal(data, &cfg)
	if parseErr != nil {
		return IConfigs{}, parseErr
	}

	return *cfg, nil
}
