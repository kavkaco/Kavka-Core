package controller

import (
	"Kavka/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService}
}

type LoginDTO struct {
	Phone string `json:"phone" xml:"phone" form:"phone"`
}

// FIXME
func (ctrl *UserController) HandleLogin(ctx *fiber.Ctx) error {
	body := new(LoginDTO)

	if err := ctx.BodyParser(body); err != nil {
		return err
	}

	// ANCHOR
	// - Fix structure to be ready to accept a complex server design
	// - Make validations and body parsing easier & DRY
	phone := body.Phone

	otp, err := ctrl.userService.Login(phone)
	if err != nil {
		return err
	}
	return nil
}
