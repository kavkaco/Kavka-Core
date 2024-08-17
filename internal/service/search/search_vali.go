package search

type SearchValidation struct {
	Input string `validate:"required,min=3"`
}
