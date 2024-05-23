package message

import (
	"github.com/kavkaco/Kavka-Core/internal/repository"
)

type UserService interface {
}

type UserManager struct {
	messageRepo repository.MessageRepository
	chatRepo    repository.ChatRepository
}

func NewMessageService(messageRepo repository.MessageRepository, chatRepo repository.ChatRepository) UserService {
	return &UserManager{messageRepo, chatRepo}
}
