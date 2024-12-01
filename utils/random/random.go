package random

import (
	"crypto/rand"
	"encoding/binary"
	"math"
)

const UserIDLength = 8

var (
	minInt = int64(math.Pow(10, float64(UserIDLength)-1))
	maxInt = int64(math.Pow(10, float64(UserIDLength))) - 1
)

func GenerateUserID() int {
	var randomNumber int64
	buf := make([]byte, 8)

	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	randomNumber = int64(binary.BigEndian.Uint64(buf))%(maxInt-minInt+1) + minInt

	return int(randomNumber)
}
