package main

import (
	"Tahagram/internal/auth"
	"testing"
)

func TestMakeVerificCode(t *testing.T) {
	t.Logf("Random Verific Code : %d\n", auth.MakeVerificCode())
}
