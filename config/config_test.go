package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	wd, _ := os.Getwd()

	_, err := Read(wd + "/configs.yml")
	assert.NoError(t, err)

	t.Log(ENV)
}
