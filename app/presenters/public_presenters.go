package presenters

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

type SimpleMessageDto struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ResponseBadRequest(ctx *gin.Context) error {
	code := fiber.ErrBadRequest.Code

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
