package random

import (
	"encoding/hex"
	"math"
	"math/rand"
	"time"
)

const OTP_LENGTH = 6

func GenerateOTP() int {
	rand.Seed(time.Now().UnixNano())

	min := math.Pow(10, OTP_LENGTH-1)
	max := math.Pow(10, OTP_LENGTH) - 1

	return rand.Intn(int(max-min) + int(min))
}

func GenerateRandomFileName(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}

	return hex.EncodeToString(bytes)
}
