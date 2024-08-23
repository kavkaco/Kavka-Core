package user

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/utils/vali"
)

type UserService interface {
	UpdateProfile(ctx context.Context, userID model.UserID, name, lastName, username, biography string) *vali.Varror
}

type UserManager struct {
	userRepo  repository.UserRepository
	validator *vali.Vali
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &UserManager{
		userRepo:  userRepo,
		validator: vali.Validator(),
	}
}

func (s *UserManager) UpdateProfile(ctx context.Context, userID model.UserID, name, lastName, username, biography string) *vali.Varror {
	validationErrors := s.validator.Validate(UpdateProfileValidation{name, lastName, username})
	if len(validationErrors) > 0 {
		return &vali.Varror{ValidationErrors: validationErrors}
	}
	user, err := s.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		return &vali.Varror{Error: ErrNotFound}
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
		return &vali.Varror{Error: ErrUpdateUser}
	}

	return nil
}
