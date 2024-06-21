package message

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/utils/hash"
)

type UserService interface {
	UpdateProfile(ctx context.Context, userID model.UserID, name, lastName, username, biography string) error
	DeleteAccount(ctx context.Context, userId model.UserID, password string) error
}

type UserManager struct {
	userRepo    repository.UserRepository
	authRepo    repository.AuthRepository
	hashManager *hash.HashManager
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

func (s *UserManager) DeleteAccount(ctx context.Context, userId model.UserID, password string) error {
	auth, err := s.authRepo.GetUserAuth(ctx, userId)
	if err != nil {
		return ErrNotFound
	}
	if !auth.EmailVerified {
		return ErrDeleteUser
	}
	validPassword := s.hashManager.CheckPasswordHash(password, auth.PasswordHash)
	if !validPassword {
		return ErrDeleteUser
	}
	err = s.userRepo.DeleteByID(ctx, userId)
	if err != nil {
		return ErrDeleteUser
	}
	return nil
}
