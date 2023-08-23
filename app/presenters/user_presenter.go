package presenters

import (
	"Kavka/domain/user"
	"Kavka/pkg/session"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func SendTokensHeader(ctx *fiber.Ctx, tokens session.LoginTokens) {
	ctx.Response().Header.Set("Refresh", tokens.RefreshToken)
	ctx.Response().Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))
}

func ResponseUserInfo(ctx *fiber.Ctx, userInfo *user.User) {
	ctx.Status(200).JSON(struct {
		Message  string
		Code     int
		UserInfo *user.User `json:"User"`
	}{
		Message:  "Success",
		Code:     200,
		UserInfo: userInfo,
	})
}
