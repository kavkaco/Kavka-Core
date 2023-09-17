package session

import (
	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/pkg/jwt_manager"

	"github.com/go-redis/redis/v8"
)

type loginPayload struct {
	OTP int `json:"otp_code"`
}

type Session struct {
	redisClient *redis.Client
	authConfigs config.Auth
	jwtManager  *jwt_manager.JwtManager
}

type LoginTokens struct {
	AccessToken  string
	RefreshToken string
}
