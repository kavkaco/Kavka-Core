package auth

import (
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

func MakeVerificCodeExpire(now time.Time) time.Time {
	return now.Add(10 * time.Minute)
}

func GetEmailWithoutAt(email string) string {
	return email[:strings.IndexByte(email, '@')]
}

func UserLimited(limitDate *time.Time) bool {
	if limitDate == nil {
		return false
	}
	return false // TODO
}
