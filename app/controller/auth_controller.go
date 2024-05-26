package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kavkaco/Kavka-Core/app/presenters"
	validator "github.com/kavkaco/Kavka-Core/app/validator"
	auth "github.com/kavkaco/Kavka-Core/internal/service/auth"
	chat "github.com/kavkaco/Kavka-Core/internal/service/chat"
	"github.com/kavkaco/Kavka-Core/pkg/email"
	"go.uber.org/zap"
)

type AuthController struct {
	ctx          context.Context
	logger       *zap.Logger
	authService  auth.AuthService
	chatService  chat.ChatService
	emailService email.EmailOtp
}

func NewAuthController(ctx context.Context, logger *zap.Logger, authService auth.AuthService, chatService chat.ChatService, emailService email.EmailOtp) *AuthController {
	return &AuthController{ctx, logger, authService, chatService, emailService}
}

func (ctrl *AuthController) HandleLogin(c *gin.Context) {
	if body := validator.ValidateBody[validator.AuthLoginRequest](c); body != nil {
		user, accessToken, refreshToken, err := ctrl.authService.Login(ctrl.ctx, body.Email, body.Password)
		if err != nil {
			presenters.ErrorResponse(c, err)
			return
		}

		userChats, err := ctrl.chatService.GetUserChats(ctrl.ctx, user.UserID)
		if err != nil {
			presenters.ErrorResponse(c, err)
			return
		}

		c.Header(presenters.RefreshTokenHeaderName, refreshToken)
		c.Header(presenters.AccessTokenHeaderName, accessToken)

		err = presenters.UserResponse(c, user, userChats)
		if err != nil {
			presenters.InternalServerErrorResponse(c)
		}
	}
}

func (ctrl *AuthController) HandleRegister(c *gin.Context) {
	if body := validator.ValidateBody[validator.AuthRegisterRequest](c); body != nil {
		_, verifyEmailToken, err := ctrl.authService.Register(ctrl.ctx, body.Name, body.LastName, body.Username, body.Email, body.Password)
		if err != nil {
			presenters.ErrorResponse(c, err)
			return
		}

		// FIXME
		fmt.Println(verifyEmailToken)

		c.JSON(http.StatusOK, presenters.CodeMessageDto{
			Code:    200,
			Message: "account registered successfully",
		})
	}
}

func (ctrl *AuthController) HandleRefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader(presenters.RefreshTokenHeaderName)
	accessToken := c.GetHeader(presenters.AccessTokenHeaderName)

	newAccessToken, err := ctrl.authService.RefreshToken(ctrl.ctx, refreshToken, accessToken)
	if errors.Is(err, auth.ErrAccessDenied) {
		presenters.AccessDenied(c)
		return
	} else if err != nil {
		presenters.ErrorResponse(c, err)
		return
	}

	c.Header(presenters.AccessTokenHeaderName, newAccessToken)

	c.JSON(http.StatusOK, presenters.CodeMessageDto{
		Code:    200,
		Message: "access token refreshed",
	})
}

func (ctrl *AuthController) HandleAuthenticate(c *gin.Context) {
	accessToken := c.GetHeader(presenters.AccessTokenHeaderName)

	// get the user info
	user, err := ctrl.authService.Authenticate(ctrl.ctx, accessToken)
	if err != nil {
		presenters.AccessDenied(c)
		return
	}

	// gathering user chats
	userChats, err := ctrl.chatService.GetUserChats(ctrl.ctx, user.UserID)
	if err != nil {
		presenters.InternalServerErrorResponse(c)
		return
	}

	err = presenters.UserResponse(c, user, userChats)
	if err != nil {
		presenters.InternalServerErrorResponse(c)
	}
}

func (ctrl *AuthController) HandleVerifyEmail(c *gin.Context) {
	redirectTo := c.Query("redirect_to")
	verifyEmailToken := c.Param("token")

	if verifyEmailToken == "" {
		presenters.BadRequestResponse(c)
		return
	}

	err := ctrl.authService.VerifyEmail(ctrl.ctx, verifyEmailToken)
	if err != nil {
		presenters.ErrorResponse(c, err)
		return
	}

	if redirectTo != "" {
		c.Redirect(http.StatusPermanentRedirect, redirectTo)
		return
	}

	c.JSON(http.StatusOK, presenters.CodeMessageDto{
		Code:    200,
		Message: "email verified",
	})
}

func (ctrl *AuthController) HandleChangePassword(c *gin.Context) {
	accessToken := c.GetHeader(presenters.AccessTokenHeaderName)
	if body := validator.ValidateBody[validator.ChangePasswordRequest](c); body != nil {
		err := ctrl.authService.ChangePassword(ctrl.ctx, accessToken, body.OldPassword, body.NewPassword)
		if err != nil {
			presenters.ErrorResponse(c, err)
			return
		}

		c.JSON(http.StatusOK, presenters.CodeMessageDto{
			Code:    200,
			Message: "password changed",
		})
	}
}
