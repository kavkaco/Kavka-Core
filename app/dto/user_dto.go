package dto

type UserLoginDto struct {
	Phone string `json:"phone" validate:"required"`
}

type UserVerifyOTPDto struct {
	Phone string `json:"phone" validate:"required"`
	OTP   int    `json:"otp"   validate:"required"`
}
