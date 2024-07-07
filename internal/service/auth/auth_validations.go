package auth

type RegisterValidation struct {
	Name     string `validate:"required,min=3,max=40"`
	LastName string `validate:"required,min=3,max=40"`
	Username string `validate:"required,min=6,max=20"`
	Email    string `validate:"required,email"` // Email format validation
	Password string `validate:"required,min=8"`
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
	AccessToken string `validate:"required"`
	OldPassword string `validate:"required"`
	NewPassword string `validate:"required,min=8"`
}

type RefreshTokenValidation struct {
	RefreshToken string `validate:"required"`
	AccessToken  string `validate:"required"`
}

type SendResetPasswordValidation struct {
	Email string `validate:"required,email"`
}

type SubmitResetPasswordValidation struct {
	ResetPasswordToken string `validate:"required"`
	NewPassword        string `validate:"required,min=8"`
}
