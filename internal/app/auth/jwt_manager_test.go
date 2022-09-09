package auth

import (
	"Kavka/config"
	"Kavka/internal/domain/user"
	"testing"
	"time"

	"Kavka/pkg/uuid"

	"github.com/stretchr/testify/assert"
)

func TestJWTGenerateAndVerify(t *testing.T) {
	var jwtManager = NewJwtManager(config.JWT{SecretKey: "1234", TTL: 10 * time.Second})

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
			token: "meaningless_string",
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
