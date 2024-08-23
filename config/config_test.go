package config

import (
	"errors"
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

func TestTestConfig(t *testing.T) {
	os.Setenv("KAVKA_ENV", "test")

	Read()
	require.Equal(t, CurrentEnv, Test)
}

func TestFunctionPanics(t *testing.T) {
	os.Setenv("KAVKA_ENV", "panic")

	defer func() {
		if r := recover(); r == nil {
			require.Error(t, errors.New("Expected panic but did not panic"))
		}
	}()

	Read()
}
