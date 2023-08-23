package presenters

import (
	"Kavka/pkg/session"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func SendTokensHeader(ctx *fiber.Ctx, tokens session.LoginTokens) {
	ctx.Response().Header.Set("Refresh", tokens.RefreshToken)
	ctx.Response().Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))
}
