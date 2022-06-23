package main

import (
	"Tahagram/pkg/auth"
	"testing"
)

func TestMakeVerificCode(t *testing.T) {
	t.Logf("Random Verific Code : %d\n", auth.MakeVerificCode())
	t.Logf("Email Without At : %s\n", auth.GetEmailWithoutAt("tahadostifam@mail.com"))
}
