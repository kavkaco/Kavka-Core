package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashManager(t *testing.T) {
	plainPassword := "Kavka&1234"

	hashManager := NewHashManager(DefaultHashParams)

	hashedPassword, err := hashManager.HashPassword(plainPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	t.Log(hashedPassword)

	valid := hashManager.CheckPasswordHash(plainPassword, hashedPassword)
	require.True(t, valid)

	valid = hashManager.CheckPasswordHash("invalid-plain-password", hashedPassword)
	require.False(t, valid)
}
