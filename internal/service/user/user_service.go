package message

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
)

type UserService interface {
	UpdateProfile(ctx context.Context, userID model.UserID, name, lastName, username, biography string) error
}

type UserManager struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &UserManager{
		userRepo: userRepo,
	}
}

func (s *UserManager) UpdateProfile(ctx context.Context, userID model.UserID, name, lastName, username, biography string) error {
	user, err := s.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		return ErrNotFound
	}

	if name != user.Name {
		user.Name = name
	}

	if lastName != user.LastName {
		user.LastName = lastName
	}

	if username != user.Username {
		user.Username = username
	}

	if biography != user.Biography {
		user.Biography = biography
	}

	err = s.userRepo.Update(ctx, userID, user.Name, user.LastName, user.Username, user.Biography)
	if err != nil {
		return ErrUpdateUser
	}

	return nil
}
