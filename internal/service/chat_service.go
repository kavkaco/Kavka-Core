package service

import (
	repository "Kavka/internal/repository/chat"
)

type ChatService struct {
	chatRepo *repository.ChatRepository
}

func NewChatService(chatRepo *repository.ChatRepository) *ChatService {
	return &ChatService{chatRepo}
}

func (s *UserService) NewChat(phone string) (int, error) {
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
