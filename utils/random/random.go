package random

import (
	"crypto/rand"
	"encoding/hex"
	"math"
	"math/big"
)

func generateRandomNumber(length int) int {
	min := int64(math.Pow(10, float64(length)-1))
	max := int64(math.Pow(10, float64(length))) - 1

	randomNumber, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		panic(err)
	}

	number := int(randomNumber.Int64()) + int(min)

	return number
}

func GenerateUserID() int {
	return generateRandomNumber(8)
}

func GenerateRandomFileName(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}

	return hex.EncodeToString(bytes)
}
