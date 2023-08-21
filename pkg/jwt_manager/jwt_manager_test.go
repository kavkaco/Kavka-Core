package jwt_manager

import (
	"Kavka/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var jwtManager = NewJwtManager(config.Auth{SECRET: "sample_secret", OTP_EXPIRE_SECONDS: 1 * time.Second})

const phone = "sample_phone_number"

func TestJWTGenerateAndVerifyRefreshToken(t *testing.T) {
	refreshToken, tokenErr := jwtManager.Generate(RefreshToken, phone)

	assert.Empty(t, tokenErr)
	assert.NotEmpty(t, refreshToken)

	cases := []struct {
		name  string
		token string
		err   error
	}{
		{
			name:  "valid",
			token: refreshToken,
			err:   nil,
		},
		{
			name:  "not_valid",
			token: "akdmakldmakldmaldalkdm",
			err:   ErrInvalidToken,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			_, verifyErr := jwtManager.Verify(tt.token, RefreshToken)

			assert.Equal(t, verifyErr, tt.err)
		})
	}
}
