package controller

import (
	validator "Kavka/app/validators"
	"Kavka/service"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService}
}

func (ctrl *UserController) HandleLogin(ctx *fiber.Ctx) error {
	validator.Validate[validator.UserLoginDto](ctx)

	// ctx.SendString("Phone:" + body.Phone)
	// otp, err := ctrl.userService.Login(phone)
	// if err != nil {
	// 	return err
	// }

	return nil
}
