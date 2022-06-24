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

func MakeVerificCodeExpire() time.Time {
	return time.Now().Add(configs.VerificCodeExpire)
}

func MakeVerificLimitDate() time.Time {
	return time.Now().Add(configs.VerificLimitDate)
}

func GetEmailWithoutAt(email string) string {
	return email[:strings.IndexByte(email, '@')]
}

func IsUserLimited(limitDate *time.Time) bool {
	if limitDate == nil {
		return false
	}
	return false // TODO
}

func IsVerificCodeExpired(expire time.Time) bool {
	return !expire.After(time.Now())
}
