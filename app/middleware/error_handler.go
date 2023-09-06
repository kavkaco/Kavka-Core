package middleware

import (
	"Kavka/app/presenters"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(ctx *gin.Context, err error) {
	presenters.ResponseError(ctx, err)
}
