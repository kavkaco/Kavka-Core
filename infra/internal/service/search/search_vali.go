package search

type searchValidation struct {
	Input string `validate:"required,min=3"`
}
