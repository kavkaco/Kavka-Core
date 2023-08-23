package service

import (
	"Kavka/pkg/jwt_manager"
	"Kavka/pkg/session"
	repository "Kavka/repository/user"
	"Kavka/utils/sms_otp"
	"errors"
)

type UserService struct {
	userRepo *repository.UserRepository
	session  *session.Session
	SmsOtp   *sms_otp.SMSOtp
}

func NewUserService(userRepo *repository.UserRepository, session *session.Session, smsOtp *sms_otp.SMSOtp) *UserService {
	return &UserService{userRepo, session, smsOtp}
}

func (s *UserService) Login(phone string) (int, error) {
	_, err := s.userRepo.FindOrCreateGuestUser(phone)
	if err != nil {
		return 0, err
	}

	otp, loginErr := s.session.Login(phone)
	if loginErr != nil {
		return 0, loginErr
	}

	return otp, nil
}

func (s *UserService) VerifyOTP(phone string, otp int) (*session.LoginTokens, error) {
	_, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return nil, repository.ErrUserNotFound
	}

	tokens, ok := s.session.VerifyOTP(phone, otp)
	if !ok {
		return nil, repository.ErrInvalidOtpCode
	}

	return &tokens, nil
}

// Refreshes access tokens and returns it
func (s *UserService) RefreshToken(refreshToken string, accessToken string) (string, error) {
	// Decode tokens and detect user phone
	payload, decodeRfErr := s.session.DecodeToken(refreshToken, jwt_manager.RefreshToken)
	if decodeRfErr != nil {
		return "", errors.New("invalid refresh token")
	}

	_, decodeAtErr := s.session.DecodeToken(accessToken, jwt_manager.AccessToken)
	if decodeAtErr != nil {
		return "", errors.New("invalid access token")
	}

	phone := payload.Phone

	// Generate & Refresh current access token
	newAccessToken, ok := s.session.NewAccessToken(phone)
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
