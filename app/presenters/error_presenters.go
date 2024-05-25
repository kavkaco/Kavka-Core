package presenters

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorDetail struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

type CodeMessageErrorDto struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Errors  []ErrorDetail `json:"errors"`
}

type CodeMessageDto struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func BadRequestResponse(ctx *gin.Context) error {
	code := http.StatusBadRequest

	ctx.JSON(http.StatusBadRequest, CodeMessageDto{
		Code:    code,
		Message: "bad request",
	})

	return nil
}

func ErrorResponse(ctx *gin.Context, err error) {
	code := http.StatusNotImplemented

	ctx.JSON(code, CodeMessageDto{
		Code:    code,
		Message: err.Error(),
	})
}

func InternalServerErrorResponse(ctx *gin.Context) {
	code := http.StatusInternalServerError

	ctx.JSON(code, CodeMessageDto{
		Code:    code,
		Message: "internal server error",
	})
}
