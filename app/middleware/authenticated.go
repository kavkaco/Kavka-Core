package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kavkaco/Kavka-Core/app/presenters"
	"github.com/kavkaco/Kavka-Core/internal/service"
	"github.com/kavkaco/Kavka-Core/utils/bearer"
)

func AuthenticatedMiddleware(userService service.UserService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		accessToken, bearerOk := bearer.AccessToken(ctx)

		if bearerOk {
			userInfo, err := userService.Authenticate(accessToken)
			if err != nil {
				presenters.AccessDenied(ctx)
				return
			}

			ctx.Set("user_static_id", userInfo.StaticID.Hex())

			// Process request
			ctx.Next()
		}
	}
}
