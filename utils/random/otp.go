package random

import (
	"math"
	"math/rand"
	"time"
)

const OTP_LENGTH = 6

func GenerateOTP() int {
	rand.Seed(time.Now().UnixNano())

	min := math.Pow(10, OTP_LENGTH-1)
	max := math.Pow(10, OTP_LENGTH) - 1

	return rand.Intn(int(max - min))
}
