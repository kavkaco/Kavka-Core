package random

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const OTP_LENGTH = 6

func GenerateOTP() int {
	min := int64(math.Pow(10, OTP_LENGTH-1))
	max := int64(math.Pow(10, OTP_LENGTH)) - 1

	randomNumber, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		panic(err)
	}

	number := int(randomNumber.Int64()) + int(min)

	if len(strconv.Itoa(number)) != OTP_LENGTH {
		number = GenerateOTP()
	}

	return number
}

func GenerateUsername() string {
	return fmt.Sprintf("guest_%d", GenerateOTP())
}

func GenerateRandomFileName(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}

	return hex.EncodeToString(bytes)
}
