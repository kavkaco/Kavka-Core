package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigsDirPath(t *testing.T) {
	t.Log(ConfigsDirPath())
	t.Log(ProjectRootPath)
}

func TestRead(t *testing.T) {
	configs := Read()

	require.NotEmpty(t, configs)

	t.Log(configs)
}
