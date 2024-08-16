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

	valid := hashManager.CheckPasswordHash(plainPassword, hashedPassword)
	require.True(t, valid)

	valid = hashManager.CheckPasswordHash("invalid-plain-password", hashedPassword)
	require.False(t, valid)
}

func Benchmark(b *testing.B) {
	hashManager := NewHashManager(DefaultHashParams)

	for i := 0; i < b.N; i++ {
		password := "Kavka&1234"
		hash, err := hashManager.HashPassword(password)
		if err != nil {
			b.Fatal(err)
		}

		hashManager.CheckPasswordHash(password, hash)
	}
}
