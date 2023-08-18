package config

import (
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	wd, _ := os.Getwd()

	_, err := Read(wd + "/configs.yml")
	if err != nil {
		t.Error(err)
	}
}
