package session

import (
	"github.com/go-redis/redis/v8"
	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/pkg/jwt_manager"
)

type sessionTokenData struct {
	TokenType string `json:"tokenType"`
}

type loginPayload struct {
	OTP int `json:"otpCode"`
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
