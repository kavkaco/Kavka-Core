package validator

import (
	fvldtor "github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = fvldtor.New()

func ValidateStruct(i interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(i)
	if err != nil {
		for _, err := range err.(fvldtor.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.FailedField = err.StructNamespace()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
