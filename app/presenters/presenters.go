package presenters

import "github.com/gofiber/fiber/v2"

type SimpleMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ResponseBadRequest(ctx *fiber.Ctx) error {
	code := fiber.ErrBadRequest.Code

	ctx.Status(code).JSON(SimpleMessage{
		Code:    code,
		Message: "Bad Request",
	})

	return nil
}
