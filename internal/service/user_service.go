package service

import (
	userRepository "github.com/kavkaco/Kavka-Core/internal/repository/user"
	auth_manager "github.com/kavkaco/Kavka-Core/pkg/auth_manager"
	"github.com/kavkaco/Kavka-Core/pkg/sms_service"
	"go.uber.org/zap"
)

type UserService interface {
	// Login(phone string) error
	// VerifyOTP(phone string, otp int) (*session.LoginTokens, error)
	// RefreshToken(refreshToken string, accessToken string) (string, error)
	// Authenticate(accessToken string) (*user.User, error)
}

type UserManager struct {
	logger      *zap.Logger
	userRepo    userRepository.UserRepository
	authManager auth_manager.AuthManager
	SmsOtp      *sms_service.SmsService
}

func NewUserService(logger *zap.Logger, userRepo userRepository.UserRepository, authManager auth_manager.AuthManager, smsOtp *sms_service.SmsService) UserService {
	return &UserManager{logger, userRepo, authManager, smsOtp}
}

// func (s *UserManager) Login(phone string) error {
// 	newUser := user.NewUser(phone)

// 	_, findErr := s.userRepo.FindByPhone(phone)
// 	if errors.Is(findErr, repository.ErrUserNotFound) {
// 		// User does not exist then it should be created!
// 		_, err := s.userRepo.Create(newUser)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	otp, loginErr := s.session.Login(phone)
// 	if loginErr != nil {
// 		return loginErr
// 	}

// 	err := s.SmsOtp.SendSMS(fmt.Sprintf("OTP Code: %d", otp), []string{phone})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *UserManager) VerifyOTP(phone string, otp int) (*session.LoginTokens, error) {
// 	foundUser, err := s.userRepo.FindByPhone(phone)
// 	if err != nil {
// 		return nil, repository.ErrUserNotFound
// 	}

// 	tokens, ok := s.authManager.GetOTP(ctx, phone, otp, foundUser.StaticID)
// 	if !ok {
// 		return nil, repository.ErrInvalidOtpCode
// 	}

// 	return &tokens, nil
// }

// func (s *UserManager) RefreshToken(refreshToken string, accessToken string) (string, error) {
// 	// Decode tokens and detect user phone
// 	payload, decodeRfErr := s.token_manager.DecodeToken(refreshToken, jwt_manager.RefreshToken)
// 	if decodeRfErr != nil {
// 		return "", errors.New("invalid refresh token")
// 	}

// 	_, decodeAtErr := s.session.DecodeToken(accessToken, jwt_manager.AccessToken)
// 	if decodeAtErr != nil {
// 		return "", errors.New("invalid access token")
// 	}

// 	foundUser, findErr := s.userRepo.FindByID(payload.StaticID)
// 	if findErr != nil {
// 		return "", findErr
// 	}

// 	newAccessToken, ok := s.session.NewAccessToken(foundUser.StaticID)
// 	if !ok {
// 		return "", errors.New("refreshing token failed")
// 	}

// 	delErr := s.session.Destroy(accessToken)
// 	if delErr != nil {
// 		return "", delErr
// 	}

// 	return newAccessToken, nil
// }

// func (s *UserManager) Authenticate(accessToken string) (*user.User, error) {
// 	payload, decodeErr := s.session.DecodeToken(accessToken, jwt_manager.AccessToken)
// 	if decodeErr != nil {
// 		return nil, errors.New("invalid access token")
// 	}

// 	foundUser, findErr := s.userRepo.FindByID(payload.StaticID)
// 	if findErr != nil {
// 		return nil, jwt_manager.ErrInvalidToken
// 	}

// 	return foundUser, nil
// }
