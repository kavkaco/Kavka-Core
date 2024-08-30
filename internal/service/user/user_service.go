package user

import (
	"context"
	"path/filepath"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/pkg/file_manager"
	"github.com/kavkaco/Kavka-Core/utils/vali"
)

const uploadDir = "tmp"

type UserService interface {
	UpdateProfile(ctx context.Context, userID model.UserID, name, lastName, username, biography string) *vali.Varror
	UpdateProfilePicture(ctx context.Context, userID model.UserID, filename string, data []byte) *vali.Varror
}

type UserManager struct {
	userRepo    repository.UserRepository
	validator   *vali.Vali
	fileManager file_manager.FileManager
}

func NewUserService(userRepo repository.UserRepository, fileManager file_manager.FileManager) UserService {
	return &UserManager{
		userRepo:    userRepo,
		validator:   vali.Validator(),
		fileManager: fileManager,
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

func (s *UserManager) UpdateProfilePicture(ctx context.Context, userID model.UserID, filename string, data []byte) *vali.Varror {
	validationErrors := s.validator.Validate(UpdateProfilePictureValidation{userID, filename})
	if len(validationErrors) > 0 {
		return &vali.Varror{ValidationErrors: validationErrors}
	}
	
	err := s.fileManager.SetFile(filename, filepath.Join(config.ProjectRootPath, uploadDir, userID))
	if err != nil {
		return &vali.Varror{Error: ErrPlacePicture}
	}

	err = s.fileManager.Write(data)
	if err != nil {
		return &vali.Varror{Error: ErrUpdatePicture}
	}

	return nil
}
