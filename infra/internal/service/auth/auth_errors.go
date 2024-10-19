package auth

import "errors"

// General errors
var (
	ErrInvalidValidation = errors.New("invalid validation")
	ErrAccessDenied      = errors.New("access denied")
	ErrEmailNotVerified  = errors.New("email not verified")
	ErrAccountLocked     = errors.New("account locked")
)

// User errors
var (
	ErrNotFound               = errors.New("user not found")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrEmailAlreadyExist      = errors.New("email already exists")
	ErrUsernameAlreadyExist   = errors.New("username already exists")
)

// Authentication errors
var (
	ErrGenerateToken    = errors.New("failed to generate token")
	ErrHashingPassword  = errors.New("failed to hash password")
	ErrCreateAuthStore  = errors.New("failed to create auth store")
	ErrCreateUser       = errors.New("failed to create user")
	ErrCreateEmailToken = errors.New("failed to create email token")
	ErrDestroyToken     = errors.New("failed to destroy token")
	ErrVerifyEmail      = errors.New("failed to verify email")
)

// Password management errors
var (
	ErrChangePassword           = errors.New("failed to change password")
	ErrClearFailedLoginAttempts = errors.New("failed to clear failed login attempts")
)

// Account management errors
var (
	ErrUnlockAccount   = errors.New("failed to unlock account")
	ErrLockAccount     = errors.New("failed to lock account")
	ErrDeleteUser      = errors.New("failed to delete user")
	ErrEmailWasNotSent = errors.New("email was not sent")
)
