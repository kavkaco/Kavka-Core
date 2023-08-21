package validator

import (
	"Kavka/app/middleware"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type ValidationErrorResponse struct {
	Success bool            `json:"success"`
	Errors  []ErrorResponse `json:"errors"`
}

type ErrorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

var validate = validator.New()

func Validate[Dto interface{}](ctx *fiber.Ctx) *Dto {
	validationErrors := []ErrorResponse{}

	body := new(Dto)

	if err := ctx.BodyParser(body); err != nil {
		middleware.ResponseBadRequest(ctx)
		return nil
	}

	errs := validate.Struct(body)

	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	if len(validationErrors) > 0 {
		ctx.JSON(ValidationErrorResponse{
			Success: false,
			Errors:  validationErrors,
		})

		return nil
	}

	return body
}
