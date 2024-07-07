package validate

import (
	"github.com/go-playground/validator/v10"
)

var vi = validator.New()

type ValidationError struct {
	Field string
	Tag   string
	Value string
}

func Validate(s interface{}) []ValidationError {
	var errors = []ValidationError{}

	if err := vi.Struct(s); err != nil {
		for _, item := range err.(validator.ValidationErrors) {
			var ve ValidationError

			ve.Field = item.Field()
			ve.Tag = item.Tag()
			ve.Value = item.Param()

			errors = append(errors, ve)
		}
	}

	return errors
}

// Transform validation errors into a Go built-in error type

type WrappedValidationError struct {
	error
	list []ValidationError
}

func NewWrappedValidationError(validationErrorsList []ValidationError) error {
	return &WrappedValidationError{list: validationErrorsList}
}

func (w *WrappedValidationError) String() string {
	return "Hello"
	// if w.error == nil {
	// 	return ""
	// }

	// finalStr, err := json.Marshal(w.list)
	// if err != nil {
	// 	return "unable to marshal instance of []ValidationError to json"
	// }

	// return string(finalStr)
}

func (w *WrappedValidationError) Cause() error {
	return w.error
}

func (w *WrappedValidationError) Error() string {
	return "invalid argument"
}
