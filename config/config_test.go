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
	os.Setenv("KAVKA_ENV", "development")

	configs := Read()
	require.NotEmpty(t, configs)
	require.Equal(t, CurrentEnv, Development)
}
