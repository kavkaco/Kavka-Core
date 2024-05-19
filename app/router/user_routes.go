package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kavkaco/Kavka-Core/app/controller"
	"github.com/kavkaco/Kavka-Core/internal/service"
	"go.uber.org/zap"
)

type UserRouter struct {
	userService service.UserService
	chatService service.ChatService
	ctrl        *controller.UserController
	router      *gin.RouterGroup
}

func NewUserRouter(logger *zap.Logger, router *gin.RouterGroup, userService service.UserService, chatService service.ChatService) *UserRouter {
	ctrl := controller.NewUserController(logger, userService, chatService)

	router.POST("/login", ctrl.HandleLogin)
	router.POST("/verify_otp", ctrl.HandleVerifyOTP)
	router.POST("/refresh_token", ctrl.HandleRefreshToken)
	router.POST("/authenticate", ctrl.HandleAuthenticate)

	return &UserRouter{userService, chatService, ctrl, router}
}
