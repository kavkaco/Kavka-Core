package bearer

import (
	"Kavka/app/presenters"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func extractTokenFromHeader(authHeader string) string {
	token := strings.Split(authHeader, "Bearer ")
	return token[1]
}

func AccessToken(ctx *fiber.Ctx) (string, bool) {
	headers := ctx.GetReqHeaders()

	bearerHeader := headers["Authorization"]
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

func RefreshToken(ctx *fiber.Ctx) (string, bool) {
	headers := ctx.GetReqHeaders()

	refreshToken := headers["Refresh"]
	if len(refreshToken) == 0 {
		presenters.ResponseBadRequest(ctx)
		return "", false
	}

	return refreshToken, true
}
