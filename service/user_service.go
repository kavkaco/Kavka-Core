package service

import (
	"Kavka/domain/user"
	"Kavka/pkg/session"
	repository "Kavka/repository/user"
	"Kavka/utils/sms_otp"
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
	_, findErr := s.userRepo.FindByPhone(phone)

	if findErr != nil {
		// Creating a new user with entered Phone
		_, createErr := s.userRepo.Create(&user.CreateUserData{
			Name:     "guest",
			LastName: "guest",
			Phone:    phone,
		})

		if createErr != nil {
			return 0, createErr
		}
	}

	otp, loginErr := s.session.Login(phone)
	if loginErr != nil {
		return 0, loginErr
	}

	return otp, nil
}
