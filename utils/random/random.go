package random

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const OtpLength = 6

func generateRandomNumber(length int) int {
	min := int64(math.Pow(10, float64(length)-1))
	max := int64(math.Pow(10, float64(length))) - 1

	randomNumber, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		panic(err)
	}

	number := int(randomNumber.Int64()) + int(min)

	if len(strconv.Itoa(number)) != length {
		number = GenerateOTP()
	}

	return number
}

func GenerateOTP() int {
	return generateRandomNumber(OtpLength)
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
