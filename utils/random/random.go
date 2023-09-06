package random

import (
	"encoding/hex"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const OTP_LENGTH = 6

func GenerateOTP() int {
	rand.Seed(time.Now().UnixNano())

	min := math.Pow(10, OTP_LENGTH-1)
	max := math.Pow(10, OTP_LENGTH) - 1

	number := rand.Intn(int(max-min) + int(min))

	if len(strconv.Itoa(number)) != OTP_LENGTH {
		number = GenerateOTP()
	}

	return number
}

func GenerateRandomFileName(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}

	return hex.EncodeToString(bytes)
}
