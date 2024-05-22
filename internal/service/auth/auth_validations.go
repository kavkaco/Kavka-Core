package service

type RegisterValidation struct {
	Name     string `validate:"required"`
	LastName string `validate:"required"`
	Username string `validate:"required"`
	Email    string `validate:"required,email"` // Email format validation
	Password string `validate:"required"`
}

type AuthenticateValidation struct {
	AccessToken string `validate:"required"`
}

type VerifyEmailValidation struct {
	VerifyEmailToken string `validate:"required"`
}

type LoginValidation struct {
	Email    string `validate:"required"`
	Password string `validate:"required"`
}

type ChangePasswordValidation struct {
	Email       string `validate:"required"`
	OldPassword string `validate:"required"`
	NewPassword string `validate:"required"`
}

type RefreshTokenValidation struct {
	RefreshToken string `validate:"required"`
	AccessToken  string `validate:"required"`
}

type SendResetPasswordVerificationValidation struct {
	Email string `validate:"required"`
}

type SubmitResetPasswordValidation struct {
	ResetPasswordToken string `validate:"required"`
	NewPassword        string `validate:"required"`
}
