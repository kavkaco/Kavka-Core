package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kavkaco/Kavka-Core/app/controller"
	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"github.com/kavkaco/Kavka-Core/internal/domain/user"
)

type UserRouter struct {
	userService user.Service
	chatService chat.Service
	ctrl        *controller.UserController
	router      *gin.RouterGroup
}

func NewUserRouter(router *gin.RouterGroup, userService user.Service, chatService chat.Service) *UserRouter {
	ctrl := controller.NewUserController(userService, chatService)

	router.POST("/login", ctrl.HandleLogin)
	router.POST("/verify_otp", ctrl.HandleVerifyOTP)
	router.POST("/refresh_token", ctrl.HandleRefreshToken)
	router.POST("/authenticate", ctrl.HandleAuthenticate)

	return &UserRouter{userService, chatService, ctrl, router}
}
