package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfigs struct {
	ListenPort int `yaml:"LISTEN_PORT"`
}

func ParseConfig(path string) (*AppConfigs, error) {
	var configs AppConfigs
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(data, &configs)

	return &configs, nil
}
