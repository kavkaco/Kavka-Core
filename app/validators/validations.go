package validator

import (
	"Kavka/app/presenters"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
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

func Validate[Dto interface{}](ctx *gin.Context) *Dto {
	validationErrors := []ErrorResponse{}

	body := new(Dto)

	bindErr := ctx.Bind(&body)

	if bindErr != nil {
		presenters.ResponseBadRequest(ctx)
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
		ctx.JSON(http.StatusBadRequest, ValidationErrorResponse{
			Success: false,
			Errors:  validationErrors,
		})

		return nil
	}

	return body
}
