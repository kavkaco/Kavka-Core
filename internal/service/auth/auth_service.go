package auth

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/pkg/email"
	"github.com/kavkaco/Kavka-Core/utils/hash"
	"github.com/kavkaco/Kavka-Core/utils/vali"
	auth_manager "github.com/tahadostifam/go-auth-manager"
)

const (
	VerifyEmailTokenExpr       = time.Minute * 5     // 5 minutes
	ResetPasswordTokenExpr     = time.Minute * 10    // 10 minutes
	AccessTokenExpr            = time.Minute * 30    // 30 minutes
	RefreshTokenExpr           = time.Hour * 24 * 14 // 2 weeks
	LockAccountDuration        = time.Second * 5
	MaximumFailedLoginAttempts = 5
)

type AuthService struct {
	authRepo     repository.AuthRepository
	userRepo     repository.UserRepository
	AuthService  auth_manager.AuthManager
	validator    *vali.Vali
	hashManager  *hash.HashManager
	emailService email.EmailService
}

func NewAuthService(authRepo repository.AuthRepository, userRepo repository.UserRepository, authManager auth_manager.AuthManager, hashManager *hash.HashManager, emailService email.EmailService) *AuthService {
	return &AuthService{authRepo, userRepo, authManager, vali.Validator(), hashManager, emailService}
}

type DetailedValidation struct {
	error
	Detail []string
}

func (a *AuthService) Register(ctx context.Context, name string, lastName string, username string, email string, password string, verifyEmailRedirectUrl string) (verifyEmailToken string, varror *vali.Varror) {
	varrors := a.validator.Validate(RegisterValidation{name, lastName, username, email, password})
	if len(varrors) > 0 {
		return "", &vali.Varror{ValidationErrors: varrors}
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

	verifyEmailToken, err = a.AuthService.GenerateToken(
		ctx, auth_manager.VerifyEmail,
		&auth_manager.TokenPayload{
			UUID:      savedUser.UserID,
			TokenType: auth_manager.VerifyEmail,
			CreatedAt: time.Now(),
		},
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

func (a *AuthService) Authenticate(ctx context.Context, accessToken string) (*model.User, *vali.Varror) {
	varrors := a.validator.Validate(AuthenticateValidation{accessToken})
	if len(varrors) > 0 {
		return nil, &vali.Varror{ValidationErrors: varrors}
	}

	tokenClaims, err := a.AuthService.DecodeAccessToken(ctx, accessToken)
	if err != nil {
		return nil, &vali.Varror{Error: ErrAccessDenied}
	}

	if len(strings.TrimSpace(tokenClaims.Payload.UUID)) == 0 {
		return nil, &vali.Varror{Error: ErrAccessDenied, ValidationErrors: varrors}
	}

	user, err := a.userRepo.FindByUserID(ctx, tokenClaims.Payload.UUID)
	if err != nil {
		return nil, &vali.Varror{Error: ErrAccessDenied, ValidationErrors: varrors}
	}

	return user, nil
}

func (a *AuthService) VerifyEmail(ctx context.Context, verifyEmailToken string) *vali.Varror {
	varrors := a.validator.Validate(VerifyEmailValidation{verifyEmailToken})
	if len(varrors) > 0 {
		return &vali.Varror{ValidationErrors: varrors}
	}

	tokenClaims, err := a.AuthService.DecodeToken(ctx, verifyEmailToken, auth_manager.VerifyEmail)
	if err != nil {
		return &vali.Varror{Error: ErrAccessDenied}
	}

	err = a.authRepo.VerifyEmail(ctx, tokenClaims.UUID)
	if err != nil {
		return &vali.Varror{Error: ErrVerifyEmail}
	}

	err = a.AuthService.DestroyToken(ctx, verifyEmailToken)
	if err != nil {
		return &vali.Varror{Error: ErrDestroyToken}
	}

	return nil
}

func (a *AuthService) Login(ctx context.Context, email string, password string) (_ *model.User, act string, rft string, varror *vali.Varror) {
	varrors := a.validator.Validate(LoginValidation{email, password})
	if len(varrors) > 0 {
		return nil, "", "", &vali.Varror{ValidationErrors: varrors}
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
	accessToken, err := a.AuthService.GenerateAccessToken(ctx, user.UserID, AccessTokenExpr)
	if err != nil {
		return nil, "", "", &vali.Varror{Error: ErrGenerateToken}
	}

	refreshToken, err := a.AuthService.GenerateRefreshToken(ctx, user.UserID, &auth_manager.RefreshTokenPayload{
		IPAddress:  "not implemented yet",
		UserAgent:  "not implemented yet",
		LoggedInAt: time.Duration(time.Now().UnixMilli()),
	}, RefreshTokenExpr)
	if err != nil {
		return nil, "", "", &vali.Varror{Error: ErrGenerateToken}
	}

	go a.authRepo.ClearFailedLoginAttempts(ctx, auth.UserID) // nolint

	return user, accessToken, refreshToken, nil
}

func (a *AuthService) ChangePassword(ctx context.Context, userID model.UserID, oldPassword string, newPassword string) *vali.Varror {
	varrors := a.validator.Validate(ChangePasswordValidation{oldPassword, newPassword})
	if len(varrors) > 0 {
		return &vali.Varror{ValidationErrors: varrors}
	}

	auth, err := a.authRepo.GetUserAuth(ctx, userID)
	if err != nil {
		return &vali.Varror{Error: ErrNotFound}
	}

	// Validate with old password
	validPassword := a.hashManager.CheckPasswordHash(oldPassword, auth.PasswordHash)
	if !validPassword {
		return &vali.Varror{Error: ErrInvalidEmailOrPassword}
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

func (a *AuthService) RefreshToken(ctx context.Context, userID model.UserID, refreshToken string) (string, *vali.Varror) {
	varrors := a.validator.Validate(RefreshTokenValidation{userID, refreshToken})
	if len(varrors) > 0 {
		return "", &vali.Varror{ValidationErrors: varrors}
	}

	// Let's check that tokens not be invalid or expired
	_, err := a.AuthService.DecodeRefreshToken(ctx, userID, refreshToken)
	if err != nil {
		return "", &vali.Varror{Error: ErrAccessDenied}
	}

	// Find auth with user_id
	_, err = a.authRepo.GetUserAuth(ctx, userID)
	if err != nil {
		return "", &vali.Varror{Error: ErrAccessDenied}
	}

	// Generate new access token
	newAccessToken, err := a.AuthService.GenerateAccessToken(ctx, userID, AccessTokenExpr)
	if err != nil {
		return "", &vali.Varror{Error: ErrGenerateToken}
	}

	return newAccessToken, nil
}

func (a *AuthService) SendResetPassword(ctx context.Context, email string, resetPasswordRedirectUrl string) (token string, timeout time.Duration, varror *vali.Varror) {
	varrors := a.validator.Validate(SendResetPasswordValidation{email})
	if len(varrors) > 0 {
		return "", 0, &vali.Varror{ValidationErrors: varrors}
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

	resetPasswordToken, err := a.AuthService.GenerateToken(ctx, auth_manager.ResetPassword, &auth_manager.TokenPayload{
		UUID:      auth.UserID,
		TokenType: auth_manager.ResetPassword,
		CreatedAt: time.Now(),
	}, ResetPasswordTokenExpr)
	if err != nil {
		return "", 0, &vali.Varror{Error: ErrGenerateToken}
	}

	err = a.emailService.SendResetPasswordEmail(email, resetPasswordRedirectUrl, user.Name, "10")
	if err != nil {
		return "", 0, &vali.Varror{Error: ErrEmailWasNotSent}
	}
	return resetPasswordToken, ResetPasswordTokenExpr, nil
}

func (a *AuthService) SubmitResetPassword(ctx context.Context, token string, newPassword string) *vali.Varror {
	varrors := a.validator.Validate(SubmitResetPasswordValidation{token, newPassword})
	if len(varrors) > 0 {
		return &vali.Varror{ValidationErrors: varrors}
	}

	tokenClaims, err := a.AuthService.DecodeToken(ctx, token, auth_manager.ResetPassword)
	if err != nil {
		return &vali.Varror{Error: ErrAccessDenied}
	}

	auth, err := a.authRepo.GetUserAuth(ctx, tokenClaims.UUID)
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

func (s *AuthService) DeleteAccount(ctx context.Context, userID model.UserID, password string) *vali.Varror {
	auth, err := s.authRepo.GetUserAuth(ctx, userID)
	if err != nil {
		return &vali.Varror{Error: ErrNotFound}
	}

	validPassword := s.hashManager.CheckPasswordHash(password, auth.PasswordHash)
	if !validPassword {
		return &vali.Varror{Error: ErrInvalidEmailOrPassword}
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
