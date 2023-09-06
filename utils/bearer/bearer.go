package bearer

import (
	"Kavka/app/presenters"
	"strings"

	"github.com/gin-gonic/gin"
)

func extractTokenFromHeader(authHeader string) string {
	token := strings.Split(authHeader, "Bearer ")
	return token[1]
}

func AccessToken(ctx *gin.Context) (string, bool) {
	bearerHeader := ctx.GetHeader("Authorization")

	if len(bearerHeader) == 0 {
		presenters.ResponseBadRequest(ctx)
		return "", false
	}

	accessToken := extractTokenFromHeader(bearerHeader)
	if len(accessToken) == 0 {
		presenters.ResponseBadRequest(ctx)
		return "", false
	}

	return accessToken, true
}

func RefreshToken(ctx *gin.Context) (string, bool) {
	refreshToken := ctx.GetHeader("refresh")

	if len(refreshToken) == 0 {
		presenters.ResponseBadRequest(ctx)
		return "", false
	}

	return refreshToken, true
}
