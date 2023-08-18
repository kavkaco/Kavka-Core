package jwt_manager

import (
	"Kavka/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JWT_ALGORITHM = jwt.SigningMethodHS256

// define errors
var (
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidTokenPayload     = errors.New("invalid token payload")
	ErrUnexpectedSigningMethod = errors.New("unexpected token signing method")
)

type UserPayload struct {
	StaticID string
	jwt.StandardClaims
}

type JwtManager struct {
	secretKey string
	ttl       time.Duration
}

type IJwtManager interface {
	Verify(token string) (*UserPayload, error)
	Generate(staticID string) (string, error)
}

func NewJwtManager(config config.Auth) IJwtManager {
	return &JwtManager{
		secretKey: config.JWTSecretKey,
		ttl:       config.OTP_EXPIRE_MINUTE,
	}
}

func (m *JwtManager) Generate(staticID string) (string, error) {
	payload := UserPayload{
		StaticID: staticID,
	}

	token := jwt.NewWithClaims(JWT_ALGORITHM, payload)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JwtManager) Verify(accessToken string) (*UserPayload, error) {
	token, err := jwt.ParseWithClaims(accessToken, &UserPayload{},
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

	payload, ok := token.Claims.(*UserPayload)
	if !ok {
		return nil, ErrInvalidTokenPayload
	}

	return payload, nil
}
