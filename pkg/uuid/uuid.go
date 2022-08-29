package uuid

import (
	"encoding/hex"
)

var (
	defaultLength = 8
)

func Random() string {
	bytes := make([]byte, defaultLength)

	return hex.EncodeToString(bytes)
}
