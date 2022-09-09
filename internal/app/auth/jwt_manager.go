package auth

import (
	"Kavka/config"
	"Kavka/internal/domain/user"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// define errors
var (
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidTokenClaims      = errors.New("invalid token claims")
	ErrUnexpectedSigningMethod = errors.New("unexpected token signing method")
)

type UserClaims struct {
	StaticID string
	jwt.StandardClaims
}

type JwtManager struct {
	secretKey string
	ttl       time.Duration
}

type IJwtManager interface {
	GenerateAccessToken(u *user.User) (string, error)
	VerifyAccessToken(accessToken string) (*UserClaims, error)
}

func NewJwtManager(config config.JWT) IJwtManager {
	return &JwtManager{
		secretKey: config.SecretKey,
		ttl:       config.TTL,
	}
}

func (m *JwtManager) GenerateAccessToken(u *user.User) (string, error) {
	claims := UserClaims{
		StaticID: u.StaticID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JwtManager) VerifyAccessToken(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnexpectedSigningMethod
			}

			return []byte(m.secretKey), nil
		},
	)

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, ErrInvalidTokenClaims
	}

	return claims, nil
}
