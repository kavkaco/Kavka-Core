package jwt_manager

import (
	"Kavka/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	RF_EXPIRE_DAY = 60 * 24 * time.Hour // 60 days
	AT_EXPIRE_DAY = 24 * time.Hour      // 1 days
)

var JWT_ALGORITHM = jwt.SigningMethodHS256

// define errors
var (
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidTokenType        = errors.New("invalid token type")
	ErrUnexpectedSigningMethod = errors.New("unexpected token signing method")
)

type JwtClaims struct {
	TokenType string
	Phone     string
	jwt.StandardClaims
}

type JwtManager struct {
	secretKey string
	ttl       time.Duration
}

type IJwtManager interface {
	Verify(token string) (*JwtClaims, error)
	Generate(phone string) (string, error)
}

const (
	RefreshToken string = "refresh"
	AccessToken  string = "access"
)

func NewJwtManager(configs config.Auth) *JwtManager {
	return &JwtManager{
		secretKey: configs.SECRET,
		ttl:       time.Duration(configs.OTP_EXPIRE_SECONDS * time.Second),
	}
}

func (m *JwtManager) Generate(tokenType string, phone string) (string, error) {
	claims := &JwtClaims{Phone: phone, TokenType: tokenType}

	token := jwt.NewWithClaims(JWT_ALGORITHM, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JwtManager) Verify(userToken string, tokenType string) (*JwtClaims, error) {
	claims := &JwtClaims{}

	token, err := jwt.ParseWithClaims(userToken, claims,
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

	if token.Valid {
		if claims.TokenType != tokenType {
			return nil, ErrInvalidTokenType
		}

		return claims, nil
	}

	return nil, ErrInvalidToken
}
