package middleware

import (
	"github.com/kavkaco/Kavka-Core/app/presenters"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(ctx *gin.Context, err error) {
	presenters.ResponseError(ctx, err)
}
