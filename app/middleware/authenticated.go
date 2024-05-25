package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/kavkaco/Kavka-Core/app/presenters"
	"github.com/kavkaco/Kavka-Core/internal/service/auth"
)

func AuthenticatedMiddleware(ctx context.Context, authService auth.AuthService) func(ctx *gin.Context) {
	return func(ginCtx *gin.Context) {
		accessToken := ginCtx.GetHeader(presenters.AccessTokenHeaderName)

		user, err := authService.Authenticate(ctx, accessToken)
		if err != nil {
			presenters.AccessDenied(ginCtx)
			return
		}

		ginCtx.Set("userID", user.UserID)

		// Process request
		ginCtx.Next()
	}
}
