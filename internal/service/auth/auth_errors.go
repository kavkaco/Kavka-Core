package service

import "errors"

var (
	ErrInvalidOtpCode         = errors.New("invalid otp code")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrAccessDenied           = errors.New("access denied")
	ErrEmailNotVerified       = errors.New("email not verified")
	ErrAccountLocked          = errors.New("account locked")
	ErrVerifyEmail            = errors.New("failed verify email")
	ErrHashingPassword        = errors.New("failed to hash password")
	ErrCreateAuthStore        = errors.New("failed to to create auth store")
	ErrCreateUser             = errors.New("failed to to create user")
	ErrInvalidValidation      = errors.New("failed to validate arguments")
	ErrCreateEmailToken       = errors.New("failed to create email token")
)
