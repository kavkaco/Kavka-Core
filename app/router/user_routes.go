package router

import (
	"Kavka/app/controller"
	"Kavka/service"

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

	return &UserRouter{service, ctrl, &router}
}
