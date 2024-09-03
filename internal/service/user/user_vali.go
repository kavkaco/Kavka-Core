package user

type UpdateProfileValidation struct {
	Name     string `validate:"required,min=3,max=40"`
	LastName string `validate:"required,min=3,max=40"`
	Username string `validate:"required,min=4,max=20"`
}

type UpdateProfilePictureValidation struct {
	UserID   string `validate:"required"`
	FileName string `validate:"required"`
}
