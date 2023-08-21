package middleware

import (
	"github.com/gofiber/fiber/v2"
)

type BadRequestError struct {
	StatusCode   int    `json:"code"`
	ErrorMessage string `json:"error"`
}

func ResponseBadRequest(ctx *fiber.Ctx) error {
	code := fiber.ErrBadRequest.Code

	ctx.Status(code).JSON(BadRequestError{
		StatusCode:   code,
		ErrorMessage: "Bad Request",
	})

	return nil
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.ErrBadRequest.Code

	ctx.Status(code).JSON(BadRequestError{
		StatusCode:   code,
		ErrorMessage: err.Error(),
	})

	return nil
}
