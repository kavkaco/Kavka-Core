package jwt_manager

import (
	"errors"
	"time"

	"Kavka/config"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	RF_EXPIRE_DAY = 60 * 24 * time.Hour // 60 days
	AT_EXPIRE_DAY = 24 * time.Hour      // 1 days
)

var JWT_ALGORITHM = jwt.SigningMethodHS512

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

func NewJwtManager(configs config.Auth) *JwtManager {
	return &JwtManager{
		secretKey: configs.SECRET,
		ttl:       configs.OTP_EXPIRE_SECONDS * time.Second, //nolint
	}
}

func (m *JwtManager) Generate(tokenType string, staticID primitive.ObjectID) (string, error) {
	createdAt := time.Now()
	claims := &JwtClaims{StaticID: staticID, TokenType: tokenType, CreatedAt: createdAt}

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
