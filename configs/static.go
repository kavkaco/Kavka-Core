package configs

import "time"

var (
	MaxVerificTryCount = 5
	VerificCodeExpire  = 10 * time.Minute
	VerificLimitDate   = 2 * time.Hour
)
