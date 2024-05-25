package random

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateOTP(t *testing.T) {
	otp := GenerateOTP()
	require.NotZero(t, otp)

	t.Log(otp)
}

func TestGenerateUsername(t *testing.T) {
	username := GenerateUsername()

	t.Log(username)
}

func TestGenerateUserID(t *testing.T) {
	userID := GenerateUserID()

	t.Log(userID)
}
