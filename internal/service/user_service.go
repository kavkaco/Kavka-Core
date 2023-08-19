package service

import (
	"Kavka/internal/domain/user"
	"Kavka/internal/modules/session"
	repository "Kavka/internal/repository/user"
)

type UserService struct {
	userRepo *repository.UserRepository
	session  *session.Session
}

func NewUserService(userRepo *repository.UserRepository, session *session.Session) *UserService {
	return &UserService{userRepo, session}
}

func (s *UserService) Login(phone string) (int, error) {
	newUser, findErr := s.userRepo.FindByPhone(phone)
	if findErr != nil {
		return 0, findErr
	}

	if newUser == nil {
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
