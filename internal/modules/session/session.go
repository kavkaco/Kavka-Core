package session

import (
	"Kavka/config"
	"Kavka/internal/modules/jwt_manager"
	"Kavka/internal/modules/random"
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type ISession interface {
	Login(phone string) (int, error)
	VerifyOTP(phone string, otp string) bool
	Logout(staticID string) error
	SaveToken(token string, payload jwt_manager.UserPayload) error
	DestroyToken(token string) error
}

type loginPayload struct {
	OTP int `json:"otp_code"`
}

type Session struct {
	redisClient *redis.Client
	authConfigs config.Auth
	jwtManager  jwt_manager.IJwtManager
}

func NewSession(redisClient *redis.Client, authConfigs config.Auth) *Session {
	jwtManager := jwt_manager.NewJwtManager(authConfigs)
	return &Session{redisClient, authConfigs, jwtManager}
}

func (session *Session) SaveToken(token string, userPayload jwt_manager.UserPayload) error {
	payload, _ := json.Marshal(userPayload)

	err := session.redisClient.Set(context.Background(), token, payload, session.authConfigs.OTP_EXPIRE_MINUTE).Err()
	if err != nil {
		return err
	}

	return nil
}

func (session *Session) DestroyToken(token string) error {
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

	err := session.redisClient.Set(context.Background(), phone, payload, session.authConfigs.OTP_EXPIRE_MINUTE).Err()
	if err != nil {
		return 0, err
	}

	return otp, nil
}

// VerifyOTP function is used to compare the stored otp code and the entered otp code by the user
// and then returns a boolean values thats gonna tell that otp is valid or not.
func (session *Session) VerifyOTP(phone string, otp int) bool {
	payload, getErr := session.redisClient.Get(context.Background(), phone).Result()
	if getErr != nil {
		return false
	}

	var data loginPayload
	unmarshalErr := json.Unmarshal([]byte(payload), &data)
	if unmarshalErr != nil {
		return false
	}

	if otp == data.OTP {
		return true
	}

	return false
}
