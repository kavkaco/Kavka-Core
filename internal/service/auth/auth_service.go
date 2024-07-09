package auth

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	auth_manager "github.com/kavkaco/Kavka-Core/pkg/auth_manager"
	"github.com/kavkaco/Kavka-Core/pkg/email"
	"github.com/kavkaco/Kavka-Core/utils/hash"
	"github.com/kavkaco/Kavka-Core/utils/vali"
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
	Register(ctx context.Context, name string, lastName string, username string, email string, password string, verifyEmailRedirectUrl string) (verifyEmailToken string, varror *vali.Varror)
	Authenticate(ctx context.Context, accessToken string) (*model.User, *vali.Varror)
	VerifyEmail(ctx context.Context, verifyEmailToken string) *vali.Varror
	Login(ctx context.Context, email string, password string) (_ *model.User, act string, rft string, varror *vali.Varror)
	ChangePassword(ctx context.Context, userID model.UserID, oldPassword string, newPassword string) *vali.Varror
	RefreshToken(ctx context.Context, refreshToken string, accessToken string) (string, *vali.Varror)
	SendResetPassword(ctx context.Context, email string, resetPasswordRedirectUrl string) (token string, timeout time.Duration, varror *vali.Varror)
	SubmitResetPassword(ctx context.Context, token string, newPassword string) *vali.Varror
	DeleteAccount(ctx context.Context, userID model.UserID, password string) *vali.Varror
}

type AuthManager struct {
	authRepo     repository.AuthRepository
	userRepo     repository.UserRepository
	authManager  auth_manager.AuthManager
	validator    *vali.Vali
	hashManager  *hash.HashManager
	emailService email.EmailService
}

func NewAuthService(authRepo repository.AuthRepository, userRepo repository.UserRepository, authManager auth_manager.AuthManager, hashManager *hash.HashManager, emailService email.EmailService) AuthService {
	return &AuthManager{authRepo, userRepo, authManager, vali.Validator(), hashManager, emailService}
}

type DetailedValidation struct {
	error
	Detail []string
}

func (a *AuthManager) Register(ctx context.Context, name string, lastName string, username string, email string, password string, verifyEmailRedirectUrl string) (verifyEmailToken string, varror *vali.Varror) {
	validationErrors := a.validator.Validate(RegisterValidation{name, lastName, username, email, password})
	if len(validationErrors) > 0 {
		return "", &vali.Varror{ValidationErrors: validationErrors}
	}

	// Check uniqueness of indexes
	isUnique, unUniqueFields := a.userRepo.IsIndexesUnique(ctx, email, username)
	if !isUnique {
		if slices.Contains(unUniqueFields, "email") {
			return "", &vali.Varror{Error: ErrEmailAlreadyExist}
		}

		if slices.Contains(unUniqueFields, "username") {
			return "", &vali.Varror{Error: ErrUsernameAlreadyExist}
		}

		return "", &vali.Varror{Error: repository.ErrUniqueConstraint}
	}

	userModel := model.NewUser(name, lastName, email, username)
	savedUser, err := a.userRepo.Create(ctx, userModel)
	if err != nil {
		return "", &vali.Varror{Error: ErrCreateUser}
	}

	passwordHash, err := a.hashManager.HashPassword(password)
	if err != nil {
		return "", &vali.Varror{Error: ErrHashingPassword}
	}

	authModel := model.NewAuth(savedUser.UserID, passwordHash)
	_, err = a.authRepo.Create(ctx, authModel)
	if err != nil {
		return "", &vali.Varror{Error: ErrCreateAuthStore}
	}

	verifyEmailToken, err = a.authManager.GenerateToken(
		ctx, auth_manager.VerifyEmail,
		auth_manager.NewTokenClaims(savedUser.UserID, auth_manager.VerifyEmail),
		VerifyEmailTokenExpr,
	)
	if err != nil {
		return "", &vali.Varror{Error: ErrCreateEmailToken}
	}

	err = a.emailService.SendVerificationEmail(email, verifyEmailRedirectUrl, verifyEmailToken)
	if err != nil {
		return "", &vali.Varror{Error: ErrEmailWasNotSent}
	}

	return verifyEmailToken, nil
}

func (a *AuthManager) Authenticate(ctx context.Context, accessToken string) (*model.User, *vali.Varror) {
	validationErrors := a.validator.Validate(AuthenticateValidation{accessToken})
	if len(validationErrors) > 0 {
		return nil, &vali.Varror{ValidationErrors: validationErrors}
	}

	tokenClaims, err := a.authManager.DecodeToken(ctx, accessToken, auth_manager.AccessToken)
	if err != nil {
		return nil, &vali.Varror{Error: ErrAccessDenied, ValidationErrors: validationErrors}
	}

	if len(strings.TrimSpace(tokenClaims.UserID)) == 0 {
		return nil, &vali.Varror{Error: ErrAccessDenied, ValidationErrors: validationErrors}
	}

	user, err := a.userRepo.FindByUserID(ctx, tokenClaims.UserID)
	if err != nil {
		return nil, &vali.Varror{Error: ErrAccessDenied, ValidationErrors: validationErrors}
	}

	return user, nil
}

func (a *AuthManager) VerifyEmail(ctx context.Context, verifyEmailToken string) *vali.Varror {
	validationErrors := a.validator.Validate(VerifyEmailValidation{verifyEmailToken})
	if len(validationErrors) > 0 {
		return &vali.Varror{ValidationErrors: validationErrors}
	}

	tokenClaims, err := a.authManager.DecodeToken(ctx, verifyEmailToken, auth_manager.VerifyEmail)
	if err != nil {
		return &vali.Varror{Error: ErrAccessDenied}
	}

	err = a.authRepo.VerifyEmail(ctx, tokenClaims.UserID)
	if err != nil {
		return &vali.Varror{Error: ErrVerifyEmail}
	}

	err = a.authManager.Destroy(ctx, verifyEmailToken)
	if err != nil {
		return &vali.Varror{Error: ErrDestroyToken}
	}

	return nil
}

func (a *AuthManager) Login(ctx context.Context, email string, password string) (_ *model.User, act string, rft string, varror *vali.Varror) {
	validationErrors := a.validator.Validate(LoginValidation{email, password})
	if len(validationErrors) > 0 {
		return nil, "", "", &vali.Varror{ValidationErrors: validationErrors}
	}

	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", "", &vali.Varror{Error: ErrInvalidEmailOrPassword}
	}

	auth, err := a.authRepo.GetUserAuth(ctx, user.UserID)
	if err != nil {
		return nil, "", "", &vali.Varror{Error: ErrInvalidEmailOrPassword}
	}

	if !auth.EmailVerified {
		return nil, "", "", &vali.Varror{Error: ErrEmailNotVerified}
	}

	// Check the expiration of account locked time
	if auth.AccountLockedUntil != 0 { //nolint
		now := time.Now()
		lockTime := time.Unix(auth.AccountLockedUntil, 0)

		// End of account lock!
		if now.After(lockTime) {
			err = a.authRepo.UnlockAccount(ctx, auth.UserID)
			if err != nil {
				return nil, "", "", &vali.Varror{Error: ErrUnlockAccount}
			}

			err = a.authRepo.ClearFailedLoginAttempts(ctx, auth.UserID)
			if err != nil {
				return nil, "", "", &vali.Varror{Error: ErrClearFailedLoginAttempts}
			}

			auth.AccountLockedUntil = 0
		}
	}

	// Account is still locked
	if auth.AccountLockedUntil != 0 {
		lockTime := time.Unix(auth.AccountLockedUntil, 0)
		return nil, "", "", &vali.Varror{Error: fmt.Errorf("%w until %v", ErrAccountLocked, lockTime)}
	}

	if auth.FailedLoginAttempts+1 == MaximumFailedLoginAttempts {
		err = a.authRepo.LockAccount(ctx, auth.UserID, LockAccountDuration)
		if err != nil {
			return nil, "", "", &vali.Varror{Error: ErrLockAccount}
		}
	}

	validPassword := a.hashManager.CheckPasswordHash(password, auth.PasswordHash)
	if !validPassword {
		// Increment the filed login attempts
		err = a.authRepo.IncrementFailedLoginAttempts(ctx, user.UserID)
		if err != nil {
			return nil, "", "", &vali.Varror{Error: ErrInvalidEmailOrPassword}
		}

		return nil, "", "", &vali.Varror{Error: ErrInvalidEmailOrPassword}
	}

	// Generate refresh token and access token
	accessToken, err := a.authManager.GenerateToken(ctx, auth_manager.AccessToken, auth_manager.NewTokenClaims(user.UserID, auth_manager.AccessToken), AccessTokenExpr)
	if err != nil {
		return nil, "", "", &vali.Varror{Error: ErrGenerateToken}
	}

	refreshToken, err := a.authManager.GenerateToken(ctx, auth_manager.RefreshToken, auth_manager.NewTokenClaims(user.UserID, auth_manager.RefreshToken), RefreshTokenExpr)
	if err != nil {
		return nil, "", "", &vali.Varror{Error: ErrGenerateToken}
	}

	err = a.authRepo.ClearFailedLoginAttempts(ctx, auth.UserID)
	if err != nil {
		return nil, "", "", &vali.Varror{Error: ErrClearFailedLoginAttempts}
	}

	return user, accessToken, refreshToken, nil
}

func (a *AuthManager) ChangePassword(ctx context.Context, userID model.UserID, oldPassword string, newPassword string) *vali.Varror {
	validationErrors := a.validator.Validate(ChangePasswordValidation{oldPassword, newPassword})
	if len(validationErrors) > 0 {
		return &vali.Varror{ValidationErrors: validationErrors}
	}

	auth, err := a.authRepo.GetUserAuth(ctx, userID)
	if err != nil {
		return &vali.Varror{Error: ErrNotFound}
	}

	// Validate with old password
	validPassword := a.hashManager.CheckPasswordHash(oldPassword, auth.PasswordHash)
	if !validPassword {
		return &vali.Varror{Error: ErrInvalidPassword}
	}

	newPasswordHash, err := a.hashManager.HashPassword(newPassword)
	if err != nil {
		return &vali.Varror{Error: ErrHashingPassword}
	}

	err = a.authRepo.ChangePassword(ctx, userID, newPasswordHash)
	if err != nil {
		return &vali.Varror{Error: ErrChangePassword}
	}

	return nil
}

func (a *AuthManager) RefreshToken(ctx context.Context, refreshToken string, accessToken string) (string, *vali.Varror) {
	validationErrors := a.validator.Validate(RefreshTokenValidation{refreshToken, accessToken})
	if len(validationErrors) > 0 {
		return "", &vali.Varror{ValidationErrors: validationErrors}
	}

	// Let's check that tokens not be invalid or expired
	rftClaims, err := a.authManager.DecodeToken(ctx, refreshToken, auth_manager.RefreshToken)
	if err != nil {
		return "", &vali.Varror{Error: ErrAccessDenied}
	}

	_, err = a.authManager.DecodeToken(ctx, accessToken, auth_manager.AccessToken)
	if err != nil {
		return "", &vali.Varror{Error: ErrAccessDenied}
	}

	// Find auth with user_id
	_, err = a.authRepo.GetUserAuth(ctx, rftClaims.UserID)
	if err != nil {
		return "", &vali.Varror{Error: ErrAccessDenied}
	}

	// Generate new access token
	newAccessToken, err := a.authManager.GenerateToken(ctx, auth_manager.AccessToken, auth_manager.NewTokenClaims(rftClaims.UserID, auth_manager.AccessToken), AccessTokenExpr)
	if err != nil {
		return "", &vali.Varror{Error: ErrGenerateToken}
	}

	// Expire old access token
	err = a.authManager.Destroy(ctx, accessToken)
	if err != nil {
		return "", &vali.Varror{Error: ErrDestroyToken}
	}

	return newAccessToken, nil
}

func (a *AuthManager) SendResetPassword(ctx context.Context, email string, resetPasswordRedirectUrl string) (token string, timeout time.Duration, varror *vali.Varror) {
	validationErrors := a.validator.Validate(SendResetPasswordValidation{email})
	if len(validationErrors) > 0 {
		return "", 0, &vali.Varror{ValidationErrors: validationErrors}
	}

	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", 0, &vali.Varror{Error: ErrNotFound}
	}

	auth, err := a.authRepo.GetUserAuth(ctx, user.UserID)
	if err != nil {
		return "", 0, &vali.Varror{Error: ErrNotFound}
	}

	if !auth.EmailVerified {
		return "", 0, &vali.Varror{Error: ErrEmailNotVerified}
	}

	if auth.FailedLoginAttempts >= MaximumFailedLoginAttempts {
		return "", 0, &vali.Varror{Error: fmt.Errorf("%w until: %v", ErrAccountLocked, auth.AccountLockedUntil)}
	}

	resetPasswordToken, err := a.authManager.GenerateToken(ctx, auth_manager.ResetPassword, auth_manager.NewTokenClaims(auth.UserID, auth_manager.ResetPassword), ResetPasswordTokenExpr)
	if err != nil {
		return "", 0, &vali.Varror{Error: ErrGenerateToken}
	}

	err = a.emailService.SendResetPasswordEmail(email, resetPasswordRedirectUrl, user.Name, "10")
	if err != nil {
		return "", 0, &vali.Varror{Error: ErrEmailWasNotSent}
	}
	return resetPasswordToken, ResetPasswordTokenExpr, nil
}

func (a *AuthManager) SubmitResetPassword(ctx context.Context, token string, newPassword string) *vali.Varror {
	validationErrors := a.validator.Validate(SubmitResetPasswordValidation{token, newPassword})
	if len(validationErrors) > 0 {
		return &vali.Varror{ValidationErrors: validationErrors}
	}

	tokenClaims, err := a.authManager.DecodeToken(ctx, token, auth_manager.ResetPassword)
	if err != nil {
		return &vali.Varror{Error: ErrAccessDenied}
	}

	auth, err := a.authRepo.GetUserAuth(ctx, tokenClaims.UserID)
	if err != nil {
		return &vali.Varror{Error: ErrAccessDenied}
	}

	newPasswordHash, err := a.hashManager.HashPassword(newPassword)
	if err != nil {
		return &vali.Varror{Error: ErrHashingPassword}
	}

	err = a.authRepo.ChangePassword(ctx, auth.UserID, newPasswordHash)
	if err != nil {
		return &vali.Varror{Error: ErrChangePassword}
	}

	return nil
}

func (s *AuthManager) DeleteAccount(ctx context.Context, userID model.UserID, password string) *vali.Varror {
	auth, err := s.authRepo.GetUserAuth(ctx, userID)
	if err != nil {
		return &vali.Varror{Error: ErrNotFound}
	}

	validPassword := s.hashManager.CheckPasswordHash(password, auth.PasswordHash)
	if !validPassword {
		return &vali.Varror{Error: ErrInvalidPassword}
	}

	err = s.authRepo.DeleteByID(ctx, userID)
	if err != nil {
		return &vali.Varror{Error: ErrDeleteUser}
	}

	err = s.userRepo.DeleteByID(ctx, userID)
	if err != nil {
		return &vali.Varror{Error: ErrDeleteUser}
	}

	return nil
}
