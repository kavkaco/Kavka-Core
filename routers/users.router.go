package routers

import (
	"Tahagram/controllers"

	"github.com/gofiber/fiber/v2"
)

func InitUsers(router fiber.Router) {
	usersRouter := router.Group("/users")

	usersRouter.Post("/signin", controllers.SigninAction)
	usersRouter.Post("/verify_code", controllers.VerifyCodeAction)
	usersRouter.Post("/authentication", controllers.AuthenticationAction)
}
