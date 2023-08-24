package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigsDirPath(t *testing.T) {
	t.Log(ConfigsDirPath())
	t.Log(ProjectRootPath)
}

func TestRead(t *testing.T) {
	configs := Read()

	assert.NotEmpty(t, configs)

	t.Log(configs)
}
