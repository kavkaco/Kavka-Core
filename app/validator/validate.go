package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kavkaco/Kavka-Core/app/presenters"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func ValidateBody[Dto interface{}](ctx *gin.Context) *Dto {
	var details []presenters.ErrorDetail

	json := new(Dto)

	if err := ctx.ShouldBindJSON(&json); err != nil {
		presenters.BadRequestResponse(ctx) //nolint
		return nil
	}

	if errs := validate.Struct(json); errs != nil {
		for _, err := range errs.(validator.ValidationErrors) { //nolint
			var ed presenters.ErrorDetail

			ed.FailedField = err.Field()
			ed.Tag = err.Tag()
			ed.Value = err.Value()
			ed.Error = true

			details = append(details, ed)
		}
	}

	if len(details) > 0 {
		ctx.JSON(http.StatusBadRequest, presenters.CodeMessageErrorDto{
			Code:   http.StatusBadRequest,
			Errors: details,
		})

		return nil
	}

	return json
}
