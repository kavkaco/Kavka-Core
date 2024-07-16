package random

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
)

const (
	OTPLength    = 6
	UserIDLength = 8
)

func generateRandomNumber(length int) int {
	min := int64(math.Pow(10, float64(length)-1))
	max := int64(math.Pow(10, float64(length))) - 1

	randomNumber, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		panic(err)
	}

	number := int(randomNumber.Int64()) + int(min)

	for lenInt(number) != length {
		number = generateRandomNumber(length)
	}

	return number
}

func GenerateOTP() int {
	return generateRandomNumber(OTPLength)
}

func GenerateUserID() int {
	return generateRandomNumber(8)
}

func GenerateUsername() string {
	return fmt.Sprintf("guest_%d", generateRandomNumber(10))
}

func GenerateRandomFileName(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}

	return hex.EncodeToString(bytes)
}

func lenInt(n int) int {
	if n == 0 {
		return 1
	}

	count := 0
	for n != 0 {
		n /= 10
		count++
	}

	return count
}
