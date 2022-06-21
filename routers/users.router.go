package routers

import (
	"Tahagram/controllers"

	"github.com/gofiber/fiber/v2"
)

func InitUsers(router fiber.Router) {
	usersRouter := router.Group("/users")

	usersRouter.Get("/login", controllers.LoginAction)
}
