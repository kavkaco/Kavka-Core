package controller

import (
	validator "Kavka/app/validators"
	"Kavka/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService}
}

func (ctrl *UserController) HandleLogin(ctx *fiber.Ctx) error {
	body := validator.Validate[validator.UserLoginDto](ctx)
	phone := body.Phone

	otp, err := ctrl.userService.Login(phone)
	if err != nil {
		return err
	}

	fmt.Printf("OTP Code: %d\n", otp)

	return nil
}
