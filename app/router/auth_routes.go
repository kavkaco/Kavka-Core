package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/kavkaco/Kavka-Core/app/controller"
	auth "github.com/kavkaco/Kavka-Core/internal/service/auth"
	chat "github.com/kavkaco/Kavka-Core/internal/service/chat"
	"github.com/kavkaco/Kavka-Core/pkg/email"
	"go.uber.org/zap"
)

func NewAuthRouter(ctx context.Context, logger *zap.Logger, router *gin.RouterGroup, authService auth.AuthService, chatService chat.ChatService, emailService email.EmailOtp) {
	ctrl := controller.NewAuthController(ctx, logger, authService, chatService, emailService)

	router.POST("/login", ctrl.HandleLogin)
	router.POST("/register", ctrl.HandleRegister)
	router.POST("/refresh_token", ctrl.HandleRefreshToken)
	router.POST("/authenticate", ctrl.HandleAuthenticate)
	router.GET("/verify_email/:token", ctrl.HandleVerifyEmail)
	router.POST("/change_password", ctrl.HandleChangePassword)
}
