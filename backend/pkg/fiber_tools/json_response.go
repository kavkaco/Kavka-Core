package fiber_tools

import "github.com/gofiber/fiber/v2"

func JSON(c *fiber.Ctx, response interface{}) error {
	err := c.JSON(response)
	c.Set("content-type", "application/json; charset=utf-8")
	return err
}
