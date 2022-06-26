package configs

import "time"

var (
	MaxVerificTryCount = 5
	VerificCodeExpire  = 10 * time.Minute // 10 * time.Minute
	VerificLimitDate   = 5 * time.Second  // 2 * time.Hour
)
