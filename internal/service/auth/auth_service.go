package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	auth_manager "github.com/kavkaco/Kavka-Core/pkg/auth_manager"
	"github.com/kavkaco/Kavka-Core/pkg/email"
	"github.com/kavkaco/Kavka-Core/utils/hash"
)

const (
	ResetPasswordTokenExpr     = time.Minute * 10    // 10 minutes
	VerifyEmailTokenExpr       = time.Minute * 5     // 5 minutes
	AccessTokenExpr            = time.Hour * 24 * 2  // 2 days
	RefreshTokenExpr           = time.Hour * 24 * 14 // 2 weeks
	LockAccountDuration        = time.Second * 5
	MaximumFailedLoginAttempts = 5
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (_ *model.User, act string, rft string, _ error)
	Register(ctx context.Context, name string, lastName string, username string, email string, password string, verifyEmailRedirectUrl string) (user *model.User, verifyEmailToken string, err error)
	VerifyEmail(ctx context.Context, verifyEmailToken string) error
	SendResetPassword(ctx context.Context, email string, resetPasswordRedirectUrl string) (token string, timeout time.Duration, _ error)
	SubmitResetPassword(ctx context.Context, token string, newPassword string) error
	ChangePassword(ctx context.Context, accessToken string, oldPassword string, newPassword string) error
	Authenticate(ctx context.Context, accessToken string) (*model.User, error)
	RefreshToken(ctx context.Context, refreshToken string, accessToken string) (newAccessToken string, err error)
	DeleteAccount(ctx context.Context, userID model.UserID, password string) error
}

type AuthManager struct {
	authRepo     repository.AuthRepository
	userRepo     repository.UserRepository
	authManager  auth_manager.AuthManager
	validator    *validator.Validate
	hashManager  *hash.HashManager
	emailService email.EmailService
}

func NewAuthService(authRepo repository.AuthRepository, userRepo repository.UserRepository, authManager auth_manager.AuthManager, hashManager *hash.HashManager, emailServic email.EmailService) AuthService {
	validator := validator.New()
	return &AuthManager{authRepo, userRepo, authManager, validator, hashManager, emailServic}
}

func (a *AuthManager) Register(ctx context.Context, name string, lastName string, username string, email string, password string, verifyEmailRedirectUrl string) (user *model.User, verifyEmailToken string, err error) {
	err = a.validator.Struct(RegisterValidation{name, lastName, username, email, password})
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	userModel := model.NewUser(name, lastName, email, username)
	savedUser, err := a.userRepo.Create(ctx, userModel)
	if errors.Is(err, repository.ErrUniqueConstraint) {
		return nil, "", repository.ErrUniqueConstraint
	} else if err != nil {
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

	err = a.emailService.SendVerificationEmail(email, verifyEmailRedirectUrl)
	if err != nil {
		return nil, "", ErrEmailWasNotSent
	}

	return savedUser, verifyEmailToken, nil
}

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

	err = a.authManager.Destroy(ctx, verifyEmailToken)
	if err != nil {
		return ErrDestroyToken
	}

	return nil
}

func (a *AuthManager) Login(ctx context.Context, email string, password string) (_ *model.User, act string, rft string, _ error) {
	err := a.validator.Struct(LoginValidation{email, password})
	if err != nil {
		return nil, "", "", fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", "", ErrInvalidEmailOrPassword
	}

	auth, err := a.authRepo.GetUserAuth(ctx, user.UserID)
	if err != nil {
		return nil, "", "", ErrInvalidEmailOrPassword
	}

	if !auth.EmailVerified {
		return nil, "", "", ErrEmailNotVerified
	}

	// Check the expiration of account locked time
	if auth.AccountLockedUntil != 0 { //nolint
		now := time.Now()
		lockTime := time.Unix(auth.AccountLockedUntil, 0)

		// End of account lock!
		if now.After(lockTime) {
			err = a.authRepo.UnlockAccount(ctx, auth.UserID)
			if err != nil {
				return nil, "", "", ErrUnlockAccount
			}

			err = a.authRepo.ClearFailedLoginAttempts(ctx, auth.UserID)
			if err != nil {
				return nil, "", "", ErrClearFailedLoginAttempts
			}

			auth.AccountLockedUntil = 0
		}
	}

	// Account is still locked
	if auth.AccountLockedUntil != 0 {
		lockTime := time.Unix(auth.AccountLockedUntil, 0)
		return nil, "", "", fmt.Errorf("%w until %v", ErrAccountLocked, lockTime)
	}

	if auth.FailedLoginAttempts+1 == MaximumFailedLoginAttempts {
		err = a.authRepo.LockAccount(ctx, auth.UserID, LockAccountDuration)
		if err != nil {
			return nil, "", "", ErrLockAccount
		}
	}

	validPassword := a.hashManager.CheckPasswordHash(password, auth.PasswordHash)
	if !validPassword {
		// Increment the filed login attempts
		err = a.authRepo.IncrementFailedLoginAttempts(ctx, user.UserID)
		if err != nil {
			return nil, "", "", ErrInvalidEmailOrPassword
		}

		return nil, "", "", ErrInvalidEmailOrPassword
	}

	// Generate refresh token and access token
	accessToken, err := a.authManager.GenerateToken(ctx, auth_manager.AccessToken, auth_manager.NewTokenClaims(user.UserID, auth_manager.AccessToken), AccessTokenExpr)
	if err != nil {
		return nil, "", "", ErrGenerateToken
	}

	refreshToken, err := a.authManager.GenerateToken(ctx, auth_manager.RefreshToken, auth_manager.NewTokenClaims(user.UserID, auth_manager.RefreshToken), RefreshTokenExpr)
	if err != nil {
		return nil, "", "", ErrGenerateToken
	}

	err = a.authRepo.ClearFailedLoginAttempts(ctx, auth.UserID)
	if err != nil {
		return nil, "", "", ErrClearFailedLoginAttempts
	}

	err = a.emailService.SendWelcomeEmail(email, user.Name)
	if err != nil {
		return nil, "", "", ErrEmailWasNotSent
	}
	return user, accessToken, refreshToken, nil
}

func (a *AuthManager) ChangePassword(ctx context.Context, accessToken string, oldPassword string, newPassword string) error {
	err := a.validator.Struct(ChangePasswordValidation{accessToken, oldPassword, newPassword})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	user, err := a.Authenticate(ctx, accessToken)
	if err != nil {
		return err
	}

	auth, err := a.authRepo.GetUserAuth(ctx, user.UserID)
	if err != nil {
		return ErrNotFound
	}

	// Validate with old password
	validPassword := a.hashManager.CheckPasswordHash(oldPassword, auth.PasswordHash)
	if !validPassword {
		return ErrInvalidPassword
	}

	newPasswordHash, err := a.hashManager.HashPassword(newPassword)
	if err != nil {
		return ErrHashingPassword
	}

	err = a.authRepo.ChangePassword(ctx, user.UserID, newPasswordHash)
	if err != nil {
		return ErrChangePassword
	}

	return nil
}

func (a *AuthManager) RefreshToken(ctx context.Context, refreshToken string, accessToken string) (string, error) {
	err := a.validator.Struct(RefreshTokenValidation{refreshToken, accessToken})
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	// Let's check that tokens not be invalid or expired
	rftClaims, err := a.authManager.DecodeToken(ctx, refreshToken, auth_manager.RefreshToken)
	if err != nil {
		return "", ErrAccessDenied
	}

	_, err = a.authManager.DecodeToken(ctx, accessToken, auth_manager.AccessToken)
	if err != nil {
		return "", ErrAccessDenied
	}

	// Find auth with user_id
	_, err = a.authRepo.GetUserAuth(ctx, rftClaims.UserID)
	if err != nil {
		return "", ErrAccessDenied
	}

	// Generate new access token
	newAccessToken, err := a.authManager.GenerateToken(ctx, auth_manager.AccessToken, auth_manager.NewTokenClaims(rftClaims.UserID, auth_manager.AccessToken), AccessTokenExpr)
	if err != nil {
		return "", ErrGenerateToken
	}

	// Expire old access token
	err = a.authManager.Destroy(ctx, accessToken)
	if err != nil {
		return "", ErrDestroyToken
	}

	return newAccessToken, nil
}

func (a *AuthManager) SendResetPassword(ctx context.Context, email string, resetPasswordRedirectUrl string) (token string, timeout time.Duration, _ error) {
	err := a.validator.Struct(SendResetPasswordValidation{email})
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", 0, err
	}

	auth, err := a.authRepo.GetUserAuth(ctx, user.UserID)
	if err != nil {
		return "", 0, err
	}

	if !auth.EmailVerified {
		return "", 0, ErrEmailNotVerified
	}

	if auth.FailedLoginAttempts >= MaximumFailedLoginAttempts {
		return "", 0, fmt.Errorf("%w until: %v", ErrAccountLocked, auth.AccountLockedUntil)
	}

	resetPasswordToken, err := a.authManager.GenerateToken(ctx, auth_manager.ResetPassword, auth_manager.NewTokenClaims(auth.UserID, auth_manager.ResetPassword), ResetPasswordTokenExpr)
	if err != nil {
		return "", 0, ErrGenerateToken
	}

	err = a.emailService.SendResetPasswordEmail(email, resetPasswordRedirectUrl, user.Name, "10")
	if err != nil {
		return "", 0, ErrEmailWasNotSent
	}
	return resetPasswordToken, ResetPasswordTokenExpr, nil
}

func (a *AuthManager) SubmitResetPassword(ctx context.Context, token string, newPassword string) error {
	err := a.validator.Struct(SubmitResetPasswordValidation{token, newPassword})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidValidation, err)
	}

	tokenClaims, err := a.authManager.DecodeToken(ctx, token, auth_manager.ResetPassword)
	if err != nil {
		return ErrAccessDenied
	}

	auth, err := a.authRepo.GetUserAuth(ctx, tokenClaims.UserID)
	if err != nil {
		return ErrAccessDenied
	}

	newPasswordHash, err := a.hashManager.HashPassword(newPassword)
	if err != nil {
		return ErrHashingPassword
	}

	err = a.authRepo.ChangePassword(ctx, auth.UserID, newPasswordHash)
	if err != nil {
		return ErrChangePassword
	}

	return nil
}

func (s *AuthManager) DeleteAccount(ctx context.Context, userID model.UserID, password string) error {
	auth, err := s.authRepo.GetUserAuth(ctx, userID)
	if err != nil {
		return ErrNotFound
	}

	validPassword := s.hashManager.CheckPasswordHash(password, auth.PasswordHash)
	if !validPassword {
		return ErrDeleteUser
	}

	err = s.authRepo.DeleteByID(ctx, userID)
	if err != nil {
		return ErrDeleteUser
	}

	err = s.userRepo.DeleteByID(ctx, userID)
	if err != nil {
		return ErrDeleteUser
	}

	return nil
}
