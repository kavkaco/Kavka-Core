package auth

import (
	"Tahagram/configs"
	"math/rand"
	"strings"
	"time"
)

func MakeVerificCode() int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	max := 999999
	min := 111111

	return min + r.Intn(max-min)
}

func MakeVerificCodeExpire() int64 {
	return time.Now().Add(configs.VerificCodeExpire).Unix()
}

func MakeVerificLimitDate() int64 {
	return time.Now().Add(configs.VerificLimitDate).Unix()
}

func IsUserLimited(limitDate int64) bool {
	now := time.Now().Unix()
	return !(now < limitDate)
}

func IsVerificCodeExpired(expire int64) bool {
	now := time.Now().Unix()
	return !(now < expire)
}

func GetEmailWithoutAt(email string) string {
	return email[:strings.IndexByte(email, '@')]
}
