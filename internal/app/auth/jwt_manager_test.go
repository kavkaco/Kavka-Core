package auth

import (
	"Kavka/config"
	"Kavka/internal/domain/user"
	"testing"
	"time"

	"Kavka/pkg/uuid"

	"github.com/stretchr/testify/assert"
)

var jwtManager = NewJwtManager(config.Auth{JWTSecretKey: "1234", AT_TTL_MINUTE: 1 * time.Second})

func TestJWTGenerateAndVerify(t *testing.T) {
	// TestGenerateAccessToken
	token, generateErr := jwtManager.GenerateAccessToken(&user.User{
		StaticID: uuid.Random(),
	})

	assert.Empty(t, generateErr)
	assert.NotEmpty(t, token)

	// TestVerifyAccessToken
	verifyTests := []struct {
		name  string
		token string
		err   error
	}{
		{
			name:  "valid",
			token: token,
			err:   nil,
		},
		{
			name:  "not_valid",
			token: "akdmakldmakldmaldalkdm",
			err:   ErrInvalidToken,
		},
	}

	for _, tt := range verifyTests {
		t.Run(tt.name, func(t *testing.T) {
			_, verifyErr := jwtManager.VerifyAccessToken(tt.token)
			assert.Equal(t, verifyErr, tt.err)
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	_, err := jwtManager.GenerateRefreshToken()

	assert.Empty(t, err)
}
