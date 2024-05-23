package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/kavkaco/Kavka-Core/app/dto"
	"github.com/kavkaco/Kavka-Core/app/presenters"
	"github.com/kavkaco/Kavka-Core/internal/service"
	"github.com/kavkaco/Kavka-Core/utils/bearer"
	"go.uber.org/zap"
)

type UserController struct {
	logger      *zap.Logger
	userService service.UserService
	chatService service.ChatService
}

func NewUserController(logger *zap.Logger, userService service.UserService, chatService service.ChatService) *UserController {
	return &UserController{logger, userService, chatService}
}

func (ctrl *UserController) HandleLogin(ginCtx *gin.Context) {
	body := dto.Validate[dto.UserLoginDto](ginCtx)
	phone := body.Phone

	err := ctrl.userService.Login(phone)
	if err != nil {
		presenters.ResponseError(ginCtx, err)
		return
	}

	ginCtx.JSON(http.StatusOK, presenters.SimpleMessageDto{
		Code:    200,
		Message: "OTP Code Sent",
	})
}

func (ctrl *UserController) HandleVerifyOTP(ginCtx *gin.Context) {
	body := dto.Validate[dto.UserVerifyOTPDto](ginCtx)

	tokens, err := ctrl.userService.VerifyOTP(body.Phone, body.OTP)
	if err != nil {
		presenters.ResponseError(ginCtx, err)
		return
	}

	ginCtx.Header("refresh", tokens)
	ginCtx.Header("authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))

	ginCtx.JSON(http.StatusOK, presenters.SimpleMessageDto{
		Code:    200,
		Message: "Logged in successfully",
	})
}

func (ctrl *UserController) HandleRefreshToken(ginCtx *gin.Context) {
	refreshToken, bearerRfOk := bearer.RefreshToken(ginCtx)

	if bearerRfOk {
		accessToken, bearerAtOk := bearer.AccessToken(ginCtx)

		if bearerAtOk {
			newAccessToken, refErr := ctrl.userService.RefreshToken(refreshToken, accessToken)
			if refErr != nil {
				presenters.ResponseError(ginCtx, refErr)
				return
			}

			ginCtx.Header("refresh", newAccessToken)
			ginCtx.Header("authorization", fmt.Sprintf("Bearer %s", refreshToken))

			ginCtx.JSON(http.StatusOK, presenters.SimpleMessageDto{
				Code:    200,
				Message: "Tokens refreshed successfully",
			})
		}
	}
}

func (ctrl *UserController) HandleAuthenticate(ginCtx *gin.Context) {
	ctx := context.TODO()
	accessToken, bearerOk := bearer.AccessToken(ginCtx)

	if bearerOk {
		// get the user info
		userInfo, err := ctrl.userService.Authenticate(accessToken)
		if err != nil {
			presenters.AccessDenied(ginCtx)
			return
		}

		// gathering user chats
		userChats, err := ctrl.chatService.GetUserChats(ctx, userInfo.StaticID)
		if err != nil {
			presenters.ResponseInternalServerError(ginCtx)
			return
		}

		err = presenters.ResponseUserInfo(ginCtx, userInfo, userChats)
		if err != nil {
			presenters.ResponseInternalServerError(ginCtx)
		}
	}
}
