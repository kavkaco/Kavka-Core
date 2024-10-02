package user

type updateProfileValidation struct {
	Name     string `validate:"required,min=3,max=40"`
	LastName string `validate:"required,min=3,max=40"`
	Username string `validate:"min=4,max=20"`
}
