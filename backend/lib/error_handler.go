package lib

import "github.com/gofiber/fiber/v2"

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	if e, ok := err.(*fiber.Error); ok {
		switch e.Code {
		case 404:
			ctx.Status(fiber.StatusNotFound).SendString("Not Found")
		}
	} else {
		ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return nil
}
