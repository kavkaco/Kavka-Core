package session

import (
	"Kavka/config"
	"Kavka/pkg/jwt_manager"
	"Kavka/utils/random"
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type ISession interface {
	Login(phone string) (int, error)
	VerifyOTP(phone string, otp string) bool
	Logout(staticID string) error
	SaveToken(token string, payload jwt_manager.JwtClaims) error
	DestroyToken(token string) error
}

func NewSession(redisClient *redis.Client, authConfigs config.Auth) *Session {
	jwtManager := jwt_manager.NewJwtManager(authConfigs)
	return &Session{redisClient, authConfigs, jwtManager}
}

// "makeExpiration" returns the expiration time for a given token type.
func makeExpiration(tokenType string) time.Duration {
	var expiration time.Duration

	if tokenType == jwt_manager.RefreshToken {
		expiration = jwt_manager.RF_EXPIRE_DAY
	}

	if tokenType == jwt_manager.AccessToken {
		expiration = jwt_manager.AT_EXPIRE_DAY
	}

	return expiration
}

func (session *Session) saveToken(token string, tokenType string) error {
	expireTime := makeExpiration(tokenType)

	err := session.redisClient.Set(context.Background(), token, nil, expireTime).Err()
	if err != nil {
		return err
	}

	return nil
}

func (session *Session) destroyToken(token string) error {
	err := session.redisClient.Del(context.Background(), token).Err()
	if err != nil {
		return err
	}

	return nil
}

// Login function is used to handle the login process for a user with the given phone number.
// It's just generate an OTP code then saves it in redis store with the key `phone`.
func (session *Session) Login(phone string) (int, error) {
	otp := random.GenerateOTP()
	payload, _ := json.Marshal(loginPayload{OTP: otp})
	expiration := session.authConfigs.OTP_EXPIRE_SECONDS * time.Second

	err := session.redisClient.Set(context.Background(), phone, payload, expiration).Err()
	if err != nil {
		return 0, err
	}

	return otp, nil
}

// VerifyOTP function is used to compare the stored otp code and the entered otp code by the user
// and then returns a boolean values thats gonna tell that otp is valid or not.
func (session *Session) VerifyOTP(phone string, otp int) (LoginTokens, bool) {
	payload, getErr := session.redisClient.Get(context.Background(), phone).Result()
	if getErr != nil {
		return LoginTokens{}, false
	}

	var data loginPayload
	unmarshalErr := json.Unmarshal([]byte(payload), &data)
	if unmarshalErr != nil {
		return LoginTokens{}, false
	}

	if otp == data.OTP {
		tokens, ok := session.RefreshToken(phone)

		if !ok {
			return LoginTokens{}, false
		}

		return tokens, true
	}

	return LoginTokens{}, false
}

func (session *Session) RefreshToken(phone string) (LoginTokens, bool) {
	// Generate Tokens
	rfToken, rfErr := session.jwtManager.Generate(jwt_manager.RefreshToken, phone)
	atToken, atErr := session.jwtManager.Generate(jwt_manager.AccessToken, phone)

	if atErr != nil && rfErr != nil {
		return LoginTokens{}, false
	}

	rfSaveErr := session.saveToken(atToken, jwt_manager.RefreshToken)
	atSaveErr := session.saveToken(atToken, jwt_manager.AccessToken)

	if atSaveErr != nil && rfSaveErr != nil {
		return LoginTokens{}, false
	}

	return LoginTokens{rfToken, atToken}, true
}
