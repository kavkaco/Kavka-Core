package middleware

import (
	"Kavka/app/presenters"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.ErrBadRequest.Code

	ctx.Status(code).JSON(presenters.SimpleMessage{
		Code:    code,
		Message: err.Error(),
	})

	return nil
}
