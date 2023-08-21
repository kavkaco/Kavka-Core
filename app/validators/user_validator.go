package validator

type UserLoginDto struct {
	Phone string `json:"phone" validate:"required,min=10,max=15"`
}
