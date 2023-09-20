package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kavkaco/Kavka-Core/app/presenters"
)

func ErrorHandler(ctx *gin.Context, err error) {
	presenters.ResponseError(ctx, err)
}
