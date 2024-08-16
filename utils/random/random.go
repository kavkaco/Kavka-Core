package random

import (
	"crypto/rand"
	"encoding/binary"
	"math"
)

const UserIDLength = 8

var min = int64(math.Pow(10, float64(UserIDLength)-1))
var max = int64(math.Pow(10, float64(UserIDLength))) - 1

func GenerateUserID() int {
	var randomNumber int64
	buf := make([]byte, 8)

	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	randomNumber = int64(binary.BigEndian.Uint64(buf))%(max-min+1) + min

	return int(randomNumber)
}
