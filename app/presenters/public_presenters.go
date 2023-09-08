package presenters

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SimpleMessageDto struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ResponseBadRequest(ctx *gin.Context) error {
	code := http.StatusBadRequest

	ctx.JSON(http.StatusBadRequest, SimpleMessageDto{
		Code:    code,
		Message: "Bad Request",
	})

	return nil
}

func ResponseError(ctx *gin.Context, err error) {
	code := http.StatusInternalServerError

	ctx.JSON(code, SimpleMessageDto{
		Code:    code,
		Message: err.Error(),
	})
}
