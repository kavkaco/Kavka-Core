package auth

type registerValidation struct {
	Name     string `validate:"required,min=1,max=40"`
	LastName string `validate:"required,min=1,max=40"`
	Username string `validate:"required,min=3,max=25"`
	Email    string `validate:"required,email"` // Email format validation
	Password string `validate:"required,min=8"`
}

type authenticateValidation struct {
	AccessToken string `validate:"required"`
}

type verifyEmailValidation struct {
	VerifyEmailToken string `validate:"required"`
}

type loginValidation struct {
	Email    string `validate:"required"`
	Password string `validate:"required"`
}

type changePasswordValidation struct {
	OldPassword string `validate:"required"`
	NewPassword string `validate:"required,min=8"`
}

type refreshTokenValidation struct {
	UserID       string `validate:"required"`
	RefreshToken string `validate:"required"`
}

type sendResetPasswordValidation struct {
	Email string `validate:"required,email"`
}

type submitResetPasswordValidation struct {
	ResetPasswordToken string `validate:"required"`
	NewPassword        string `validate:"required,min=8"`
}
