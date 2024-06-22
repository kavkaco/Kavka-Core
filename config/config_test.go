package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigsDirPath(t *testing.T) {
	t.Log(ConfigsDirPath())
	t.Log(ProjectRootPath)
}

func TestDevelopmentConfig(t *testing.T) {
	os.Setenv("ENV", "development")
	configs := Read()
	require.NotEmpty(t, configs)
	require.Equal(t, configs.Mongo.Username, "mongo")
}

func TestProductionConfig(t *testing.T) {
	os.Setenv("ENV", "production")
	configs := Read()
	require.NotEmpty(t, configs)
	require.Equal(t, configs.Mongo.Username, "amir")
}
