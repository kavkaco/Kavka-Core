package service

import "errors"

var (
	ErrInvalidPassword        = errors.New("invalid password")
	ErrInvalidOtpCode         = errors.New("invalid otp code")
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrAccessDenied           = errors.New("access denied")
	ErrEmailNotVerified       = errors.New("email not verified")
	ErrAccountLocked          = errors.New("account locked")
	ErrVerifyEmail            = errors.New("failed verify email")
	ErrGenerateToken          = errors.New("failed to generate token")
	ErrHashingPassword        = errors.New("failed to hash password")
	ErrCreateAuthStore        = errors.New("failed to to create auth store")
	ErrCreateUser             = errors.New("failed to to create user")
	ErrInvalidValidation      = errors.New("failed to validate arguments")
	ErrCreateEmailToken       = errors.New("failed to create email token")
	ErrDestroyToken           = errors.New("failed to destroy token")
	ErrChangePassword         = errors.New("failed to change password")
)
