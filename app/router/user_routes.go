package router

import (
	"Kavka/app/controller"
	"Kavka/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserRouter struct {
	service *service.UserService
	ctrl    *controller.UserController
	router  *fiber.Router
}

func NewUserRouter(router fiber.Router, service *service.UserService) *UserRouter {
	ctrl := controller.NewUserController(service)

	router.Post("/login", ctrl.HandleLogin)
	router.Post("/verify_otp", ctrl.HandleVerifyOTP)
	router.Post("/refresh_token", ctrl.HandleRefreshToken)
	router.Post("/authenticate", ctrl.HandleAuthenticate)

	return &UserRouter{service, ctrl, &router}
}
