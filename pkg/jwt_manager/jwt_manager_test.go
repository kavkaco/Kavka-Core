package jwt_manager

import (
	"testing"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtManager = NewJwtManager(config.Auth{SECRET: "sample_secret"}, DEFAULT_OTP_EXPIRE)

var StaticID = primitive.NewObjectID()

func TestJWTGenerateAndVerifyRefreshToken(t *testing.T) {
	refreshToken, tokenErr := jwtManager.Generate(RefreshToken, StaticID)

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
