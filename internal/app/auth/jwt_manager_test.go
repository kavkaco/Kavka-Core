package auth

import (
	"Kavka/config"
	"Kavka/internal/domain/user"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWTGenerateAndVerify(t *testing.T) {
	var jwtManager = NewJwtManager(config.JWT{SecretKey: "1234", TTL: 20})

	// TestGenerate
	token, generateErr := jwtManager.Generate(&user.User{
		StaticID: uuid.New(),
	})

	if generateErr != nil {
		t.Error(generateErr)
	}

	assert.NotEmpty(t, token)

	// TestVerify
	userClaims, verifyErr := jwtManager.Verify(token)
	if verifyErr != nil {
		t.Error(verifyErr)
	}

	t.Log(userClaims)
}
