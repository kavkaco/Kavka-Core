package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	auth_manager "github.com/kavkaco/Kavka-Core/pkg/auth_manager"
	"github.com/kavkaco/Kavka-Core/utils/hash"
)

const VerifyEmailTokenExpr = time.Minute * 10
const MaximumFailedLoginAttempts = 5

type AuthService interface {
	Login(ctx context.Context, email, password string) (*model.User, error)
	Register(ctx context.Context, name string, lastName string, username string, email string, password string) (user *model.User, verifyEmailToken string, err error)

	VerifyEmail(ctx context.Context, verifyEmailToken string) error

	SendResetPasswordVerification(ctx context.Context, email string) (timeout time.Time, err error)
	SubmitResetPassword(ctx context.Context, token string, newPassword string) error

	ChangePassword(ctx context.Context, oldPassword string, newPassword string) error

	Authenticate(ctx context.Context, accessToken string) (*model.User, error)
	RefreshToken(ctx context.Context, refreshToken string, accessToken string) (string, error)
}

type AuthManager struct {
	authRepo    repository.AuthRepository
	userRepo    repository.UserRepository
	authManager auth_manager.AuthManager
	validator   *validator.Validate
	hashManager *hash.HashManager
}

func NewAuthService(authRepo repository.AuthRepository, userRepo repository.UserRepository, authManager auth_manager.AuthManager, hashManager *hash.HashManager) AuthService {
	validator := validator.New()
	return &AuthManager{authRepo, userRepo, authManager, validator, hashManager}
}

func (a *AuthManager) Register(ctx context.Context, name string, lastName string, username string, email string, password string) (user *model.User, verifyEmailToken string, err error) {
	err = a.validator.Struct(RegisterValidation{name, lastName, username, email, password})
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	userModel := model.NewUser(name, lastName, email, username)
	savedUser, err := a.userRepo.Create(ctx, userModel)
	if err != nil {
		return nil, "", ErrCreateUser
	}

	passwordHash, err := a.hashManager.HashPassword(password)
	if err != nil {
		return nil, "", ErrHashingPassword
	}

	authModel := model.NewAuth(savedUser.UserID, passwordHash)
	_, err = a.authRepo.Create(ctx, authModel)
	if err != nil {
		return nil, "", ErrCreateAuthStore
	}

	verifyEmailToken, err = a.authManager.GenerateToken(
		ctx, auth_manager.VerifyEmail,
		auth_manager.NewTokenClaims(savedUser.UserID, auth_manager.VerifyEmail),
		VerifyEmailTokenExpr,
	)
	if err != nil {
		return nil, "", ErrCreateEmailToken
	}

	return savedUser, verifyEmailToken, nil
}

// Authenticate implements AuthService.
func (a *AuthManager) Authenticate(ctx context.Context, accessToken string) (*model.User, error) {
	err := a.validator.Struct(AuthenticateValidation{accessToken})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	tokenClaims, err := a.authManager.DecodeToken(ctx, accessToken, auth_manager.AccessToken)
	if err != nil {
		return nil, ErrAccessDenied
	}

	if len(strings.TrimSpace(tokenClaims.UserID)) == 0 {
		return nil, ErrAccessDenied
	}

	user, err := a.userRepo.FindByUserID(ctx, tokenClaims.UserID)
	if err != nil {
		return nil, ErrAccessDenied
	}

	return user, nil
}

// VerifyEmail implements AuthService.
func (a *AuthManager) VerifyEmail(ctx context.Context, verifyEmailToken string) error {
	err := a.validator.Struct(VerifyEmailValidation{verifyEmailToken})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	tokenClaims, err := a.authManager.DecodeToken(ctx, verifyEmailToken, auth_manager.VerifyEmail)
	if err != nil {
		return ErrAccessDenied
	}

	err = a.authRepo.VerifyEmail(ctx, tokenClaims.UserID)
	if err != nil {
		return ErrVerifyEmail
	}

	return nil
}

// Login implements AuthService.
func (a *AuthManager) Login(ctx context.Context, email string, password string) (*model.User, error) {
	err := a.validator.Struct(LoginValidation{email, password})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidEmailOrPassword
	}

	auth, err := a.authRepo.GetUserAuth(ctx, user.UserID)
	if err != nil {
		return nil, ErrInvalidEmailOrPassword
	}

	if !auth.EmailVerified {
		return nil, ErrEmailNotVerified
	}

	if auth.FailedLoginAttempts >= MaximumFailedLoginAttempts {
		return nil, fmt.Errorf("%w until: %v", ErrAccountLocked, auth.AccountLockedUntil.String())
	}

	validPassword := a.hashManager.CheckPasswordHash(password, auth.PasswordHash)
	if !validPassword {
		// Increment the filed login attempts
		err := a.authRepo.IncrementFailedLoginAttempts(ctx, user.UserID)
		if err != nil {
			return nil, ErrInvalidEmailOrPassword
		}

		return nil, ErrInvalidEmailOrPassword
	}

	return user, nil
}

// ChangePassword implements AuthService.
func (a *AuthManager) ChangePassword(ctx context.Context, oldPassword string, newPassword string) error {
	panic("unimplemented")
}

// RefreshToken implements AuthService.
func (a *AuthManager) RefreshToken(ctx context.Context, refreshToken string, accessToken string) (string, error) {
	panic("unimplemented")
}

// SendResetPasswordVerification implements AuthService.
func (a *AuthManager) SendResetPasswordVerification(ctx context.Context, email string) (timeout time.Time, err error) {
	panic("unimplemented")
}

// SubmitResetPassword implements AuthService.
func (a *AuthManager) SubmitResetPassword(ctx context.Context, token string, newPassword string) error {
	panic("unimplemented")
}
