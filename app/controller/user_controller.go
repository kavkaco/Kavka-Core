package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/kavkaco/Kavka-Core/app/dto"
	"github.com/kavkaco/Kavka-Core/app/presenters"
	"github.com/kavkaco/Kavka-Core/internal/model/chat"
	"github.com/kavkaco/Kavka-Core/internal/model/user"
	"github.com/kavkaco/Kavka-Core/pkg/session"
	"github.com/kavkaco/Kavka-Core/utils/bearer"
	"go.uber.org/zap"
)

type UserController struct {
	logger      *zap.Logger
	userService user.Service
	chatService chat.Service
}

func NewUserController(logger *zap.Logger, userService user.Service, chatService chat.Service) *UserController {
	return &UserController{logger, userService, chatService}
}

func (ctrl *UserController) HandleLogin(ctx *gin.Context) {
	body := dto.Validate[dto.UserLoginDto](ctx)
	phone := body.Phone

	err := ctrl.userService.Login(phone)
	if err != nil {
		presenters.ResponseError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, presenters.SimpleMessageDto{
		Code:    200,
		Message: "OTP Code Sent",
	})
}

func (ctrl *UserController) HandleVerifyOTP(ctx *gin.Context) {
	body := dto.Validate[dto.UserVerifyOTPDto](ctx)

	tokens, err := ctrl.userService.VerifyOTP(body.Phone, body.OTP)
	if err != nil {
		presenters.ResponseError(ctx, err)
		return
	}

	presenters.SendTokensHeader(ctx, *tokens)

	ctx.JSON(http.StatusOK, presenters.SimpleMessageDto{
		Code:    200,
		Message: "Logged in successfully",
	})
}

func (ctrl *UserController) HandleRefreshToken(ctx *gin.Context) {
	refreshToken, bearerRfOk := bearer.RefreshToken(ctx)

	if bearerRfOk {
		accessToken, bearerAtOk := bearer.AccessToken(ctx)

		if bearerAtOk {
			newAccessToken, refErr := ctrl.userService.RefreshToken(refreshToken, accessToken)
			if refErr != nil {
				presenters.ResponseError(ctx, refErr)
				return
			}

			newTokens := session.LoginTokens{AccessToken: newAccessToken, RefreshToken: refreshToken}
			presenters.SendTokensHeader(ctx, newTokens)

			ctx.JSON(http.StatusOK, presenters.SimpleMessageDto{
				Code:    200,
				Message: "Tokens refreshed successfully",
			})
		}
	}
}

func (ctrl *UserController) HandleAuthenticate(ctx *gin.Context) {
	accessToken, bearerOk := bearer.AccessToken(ctx)

	if bearerOk {
		// get the user info
		userInfo, err := ctrl.userService.Authenticate(accessToken)
		if err != nil {
			presenters.AccessDenied(ctx)
			return
		}

		// gathering user chats
		userChats, err := ctrl.chatService.GetUserChats(userInfo.StaticID)
		if err != nil {
			presenters.ResponseInternalServerError(ctx)
			return
		}

		err = presenters.ResponseUserInfo(ctx, userInfo, userChats)
		if err != nil {
			presenters.ResponseInternalServerError(ctx)
		}
	}
}
