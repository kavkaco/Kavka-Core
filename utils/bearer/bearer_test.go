package bearer

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestExtractFromHeader(t *testing.T) {
	sampleToken := "Hello_World"
	authorizationToken := "Bearer " + sampleToken

	token := ExtractFromHeader(authorizationToken)

	assert.Equal(t, token, sampleToken)
}
