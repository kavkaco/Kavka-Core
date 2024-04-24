package service

import (
	"errors"
	"fmt"

	"github.com/kavkaco/Kavka-Core/internal/domain/user"
	repository "github.com/kavkaco/Kavka-Core/internal/repository/user"
	"github.com/kavkaco/Kavka-Core/pkg/jwt_manager"
	"github.com/kavkaco/Kavka-Core/pkg/session"
	"github.com/kavkaco/Kavka-Core/pkg/sms_service"
	"go.uber.org/zap"
)

type userService struct {
	logger   *zap.Logger
	userRepo user.UserRepository
	session  *session.Session
	SmsOtp   *sms_service.SmsService
}

func NewUserService(logger *zap.Logger, userRepo user.UserRepository, session *session.Session, smsOtp *sms_service.SmsService) user.Service {
	return &userService{logger, userRepo, session, smsOtp}
}

// Login function gets user's phone and find it or created it in the database,
// then generates an otp code and stores it in redis store and returns `otp code` as int and an `error`.
func (s *userService) Login(phone string) error {
	newUser := user.NewUser(phone)

	_, findErr := s.userRepo.FindByPhone(phone)
	if errors.Is(findErr, repository.ErrUserNotFound) {
		// User does not exist then it should be created!
		_, err := s.userRepo.Create(newUser)
		if err != nil {
			return err
		}
	}

	otp, loginErr := s.session.Login(phone)
	if loginErr != nil {
		return loginErr
	}

	err := s.SmsOtp.SendSMS(fmt.Sprintf("OTP Code: %d", otp), []string{phone})
	if err != nil {
		return err
	}

	return nil
}

// VerifyOTP function gets phone and otp code and checks if the otp code was correct for
// mentioned phone, it's going to return an instance of *session.LoginTokens and an error.
func (s *userService) VerifyOTP(phone string, otp int) (*session.LoginTokens, error) {
	foundUser, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return nil, repository.ErrUserNotFound
	}

	tokens, ok := s.session.VerifyOTP(phone, otp, foundUser.StaticID)
	if !ok {
		return nil, repository.ErrInvalidOtpCode
	}

	return &tokens, nil
}

// RefreshToken function is used to refresh `Access Token`, It's returns a new `Access Token` and an error.
func (s *userService) RefreshToken(refreshToken string, accessToken string) (string, error) {
	// Decode tokens and detect user phone
	payload, decodeRfErr := s.session.DecodeToken(refreshToken, jwt_manager.RefreshToken)
	if decodeRfErr != nil {
		return "", errors.New("invalid refresh token")
	}

	_, decodeAtErr := s.session.DecodeToken(accessToken, jwt_manager.AccessToken)
	if decodeAtErr != nil {
		return "", errors.New("invalid access token")
	}

	foundUser, findErr := s.userRepo.FindByID(payload.StaticID)
	if findErr != nil {
		return "", findErr
	}

	// Generate & Refresh current access token
	newAccessToken, ok := s.session.NewAccessToken(foundUser.StaticID)
	if !ok {
		return "", errors.New("refreshing token failed")
	}

	// Destroy old token
	delErr := s.session.Destroy(accessToken)
	if delErr != nil {
		return "", delErr
	}

	return newAccessToken, nil
}

// Authenticate function is used to authenticate a user and returns a `*user.User` and an error.
func (s *userService) Authenticate(accessToken string) (*user.User, error) {
	payload, decodeErr := s.session.DecodeToken(accessToken, jwt_manager.AccessToken)
	if decodeErr != nil {
		return nil, errors.New("invalid access token")
	}

	foundUser, findErr := s.userRepo.FindByID(payload.StaticID)
	if findErr != nil {
		return nil, jwt_manager.ErrInvalidToken
	}

	return foundUser, nil
}
