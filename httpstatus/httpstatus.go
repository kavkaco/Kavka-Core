package httpstatus

import "github.com/gofiber/fiber/v2"

func InternalServerError(c *fiber.Ctx) {
	c.Status(500).JSON(fiber.Map{
		"message": "An error occurred on the server side",
	})
}
