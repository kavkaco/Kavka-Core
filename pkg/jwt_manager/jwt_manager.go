package jwt_manager

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kavkaco/Kavka-Core/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	RfExpireDay = 60 * 24 * time.Hour // 60 days
	AtExpireDay = 24 * time.Hour      // 1 day
)

var JwtAlgorithm = jwt.SigningMethodHS512

// define errors.
var (
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidTokenType        = errors.New("invalid token type")
	ErrUnexpectedSigningMethod = errors.New("unexpected token signing method")
)

type JwtClaims struct {
	TokenType string
	StaticID  primitive.ObjectID
	CreatedAt time.Time
	jwt.StandardClaims
}

type JwtManager struct {
	secretKey string
	ttl       time.Duration
}

const (
	RefreshToken string = "refresh"
	AccessToken  string = "access"
)

const DefaultOtpExpire = 120 * time.Second

func NewJwtManager(configs config.Auth, otpExpire time.Duration) *JwtManager {
	return &JwtManager{
		secretKey: configs.SECRET,
		ttl:       otpExpire,
	}
}

func (m *JwtManager) Generate(tokenType string, staticID primitive.ObjectID) (string, error) {
	createdAt := time.Now()
	claims := &JwtClaims{StaticID: staticID, TokenType: tokenType, CreatedAt: createdAt}

	token := jwt.NewWithClaims(JwtAlgorithm, claims)
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
