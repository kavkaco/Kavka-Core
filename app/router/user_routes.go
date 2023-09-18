package router

import (
	"github.com/kavkaco/Kavka-Core/app/controller"
	"github.com/kavkaco/Kavka-Core/internal/service"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	service *service.UserService
	ctrl    *controller.UserController
	router  *gin.RouterGroup
}

func NewUserRouter(router *gin.RouterGroup, service *service.UserService) *UserRouter {
	ctrl := controller.NewUserController(service)

	router.POST("/login", ctrl.HandleLogin)
	router.POST("/verify_otp", ctrl.HandleVerifyOTP)
	router.POST("/refresh_token", ctrl.HandleRefreshToken)
	router.POST("/authenticate", ctrl.HandleAuthenticate)

	return &UserRouter{service, ctrl, router}
}
