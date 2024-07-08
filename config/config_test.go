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

	devConfigs := Read()
	require.NotEmpty(t, devConfigs)
	require.Equal(t, CurrentEnv, Development)

	os.Setenv("KAVKA_ENV", "test")

	testConfigs := Read()
	require.NotEmpty(t, testConfigs)
	require.Equal(t, CurrentEnv, Test)
}
