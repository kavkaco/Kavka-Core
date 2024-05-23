package auth_manager

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/utils/random"
)

var (
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidTokenType        = errors.New("invalid token type")
	ErrUnexpectedSigningMethod = errors.New("unexpected token signing method")
	ErrNotFound                = errors.New("not found")
)

var TokenEncodingAlgorithm = jwt.SigningMethodHS512

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
	ResetPassword
	VerifyEmail
)

type AuthManager interface {
	GenerateToken(ctx context.Context, tokenType TokenType, tokenPayload *TokenClaims, expr time.Duration) (token string, err error)
	DecodeToken(ctx context.Context, token string, tokenType TokenType) (claims *TokenClaims, err error)
	Destroy(ctx context.Context, key string) (err error)
	GetOTP(ctx context.Context, uniqueID string) (otp string, err error)
	SetOTP(ctx context.Context, uniqueID string, expr time.Duration) (otp string, err error)
}

type AuthManagerOpts struct {
	RefreshTokenExpiration time.Duration
	AccessTokenExpiration  time.Duration
	PrivateKey             string
}

// Used as jwt claims
type TokenClaims struct {
	UserID    model.UserID `json:"userID"`
	CreatedAt time.Time    `json:"createdAt"`
	TokenType TokenType    `json:"tokenType"`
	jwt.StandardClaims
}

func NewTokenClaims(userID model.UserID, tokenType TokenType) *TokenClaims {
	return &TokenClaims{
		UserID:    userID,
		CreatedAt: time.Now(),
		TokenType: tokenType,
	}
}

type authManager struct {
	redisClient *redis.Client
	opts        AuthManagerOpts
}

func NewAuthManager(redisClient *redis.Client, opts AuthManagerOpts) AuthManager {
	return &authManager{redisClient, opts}
}

func (t *authManager) GenerateToken(ctx context.Context, tokenType TokenType, tokenClaims *TokenClaims, expr time.Duration) (_ string, _ error) {
	token, err := jwt.NewWithClaims(TokenEncodingAlgorithm, tokenClaims).SignedString([]byte(t.opts.PrivateKey))
	if err != nil {
		return "", err
	}

	cmd := t.redisClient.Set(ctx, token, nil, expr)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}

	return token, nil
}

func (t *authManager) DecodeToken(ctx context.Context, token string, tokenType TokenType) (_ *TokenClaims, _ error) {
	_, err := t.redisClient.Get(ctx, token).Result()
	if err != nil {
		return nil, err
	}

	tokenClaims := &TokenClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, tokenClaims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnexpectedSigningMethod
			}

			return []byte(t.opts.PrivateKey), nil
		},
	)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if jwtToken.Valid {
		if tokenClaims.TokenType != tokenType {
			return nil, ErrInvalidTokenType
		}

		return tokenClaims, nil
	}

	return &TokenClaims{}, ErrInvalidToken
}

func (t *authManager) Destroy(ctx context.Context, key string) (_ error) {
	cmd := t.redisClient.Del(ctx, key)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (t *authManager) GetOTP(ctx context.Context, uniqueID string) (_ string, _ error) {
	result, err := t.redisClient.Get(ctx, uniqueID).Result()
	if err != nil {
		return "", err
	}

	if len(strings.TrimSpace(result)) > 0 {
		return result, nil
	}

	return "", ErrNotFound
}

func (t *authManager) SetOTP(ctx context.Context, uniqueID string, expr time.Duration) (_ string, _ error) {
	otp := fmt.Sprintf("%d", random.GenerateOTP())

	_, err := t.redisClient.Set(ctx, uniqueID, otp, expr).Result()
	if err != nil {
		return "", err
	}

	return string(otp), nil
}
