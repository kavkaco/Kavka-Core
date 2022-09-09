package auth

import (
	"Kavka/internal/adapter/cache"
	"Kavka/internal/domain/user"
)

type AuthService struct {
	userCache      *cache.UserCache
	userRepository *user.Repository
}

func NewAuthService(userCache *cache.UserCache, userRepository *user.Repository) *AuthService {
	return &AuthService{userCache, userRepository}
}
