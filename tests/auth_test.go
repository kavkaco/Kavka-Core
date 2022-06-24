package main

import (
	"Tahagram/pkg/auth"
	"testing"
	"time"
)

func TestMakeVerificCode(t *testing.T) {
	t.Logf("Random Verific Code : %d\n", auth.MakeVerificCode())
}

func TestEmailWithoutAt(t *testing.T) {
	t.Logf("Email Without At : %s\n", auth.GetEmailWithoutAt("tahadostifam@mail.com"))
}

func TestISVerificCodeExpired(t *testing.T) {
	t.Logf("Expired: %v", auth.IsVerificCodeExpired(time.Now().Truncate(1*time.Minute)))

	t.Logf("Expired: %v", auth.IsVerificCodeExpired(time.Now().Add(1*time.Minute)))
}
