package service

import (
	"Kavka/internal/domain/user"
	repository "Kavka/internal/repository/user"
	"Kavka/pkg/jwt_manager"
	"Kavka/pkg/session"
	"Kavka/utils/sms_otp"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

type UserService struct {
	userRepo *repository.UserRepository
	session  *session.Session
	SmsOtp   *sms_otp.SMSOtp
}

func NewUserService(userRepo *repository.UserRepository, session *session.Session, smsOtp *sms_otp.SMSOtp) *UserService {
	return &UserService{userRepo, session, smsOtp}
}

func (s *UserService) FindByUsername(username string) (*user.User, error) {
	users, err := s.userRepo.Where(bson.D{{Key: "username", Value: username}})
	if err != nil {
		return nil, err
	}

	if len(users) > 0 && users[0] != nil {
		return users[0], nil
	}

	return nil, repository.ErrUserNotFound
}

// Login function gets user's phone and find it or created it in the database,
// then generates a otp code and stores it in redis store and returns `otp code` as int and an `error`
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

// VerifyOTP function gets phone and otp code and checks if the otp code was correct for
// mentioned phone, its gonna return an instance of *session.LoginTokens and an error.
func (s *UserService) VerifyOTP(phone string, otp int) (*session.LoginTokens, error) {
	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return nil, repository.ErrUserNotFound
	}

	tokens, ok := s.session.VerifyOTP(phone, otp, user.StaticID)
	if !ok {
		return nil, repository.ErrInvalidOtpCode
	}

	return &tokens, nil
}

// RefreshToken function is used to refresh `Access Token`, It's returns a new `Access Token` and an error.
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

	user, findErr := s.userRepo.FindByID(payload.StaticID)
	if findErr != nil {
		return "", findErr
	}

	// Generate & Refresh current access token
	newAccessToken, ok := s.session.NewAccessToken(user.StaticID)
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

// `Authenticate` function is used to authenticate a user and returns a `*user.User` and an error.
func (s *UserService) Authenticate(accessToken string) (*user.User, error) {
	payload, decodeErr := s.session.DecodeToken(accessToken, jwt_manager.AccessToken)
	if decodeErr != nil {
		return nil, errors.New("invalid access token")
	}

	user, findErr := s.userRepo.FindByID(payload.StaticID)
	if findErr != nil {
		return nil, jwt_manager.ErrInvalidToken
	}

	return user, nil
}
