package routers

import (
	"Tahagram/controllers"

	"github.com/gofiber/fiber/v2"
)

func InitUsers(app *fiber.App) {
	app.Post("/users/login", controllers.LoginAction)
}
