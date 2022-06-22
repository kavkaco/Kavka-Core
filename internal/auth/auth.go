package auth

import (
	"math/rand"
	"strings"
	"time"
)

func MakeVerificCode() int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return r.Intn(999999)
}

func GetEmailWithoutAt(email string) string {
	return email[:strings.IndexByte(email, '@')]
}
