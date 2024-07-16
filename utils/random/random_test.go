package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateOTP(t *testing.T) {
	otp := GenerateOTP()
	assert.Equal(t, OTPLength, lenInt(otp))
}

func TestGenerateUserID(t *testing.T) {
	userID := GenerateUserID()
	assert.Equal(t, UserIDLength, lenInt(userID))
}
