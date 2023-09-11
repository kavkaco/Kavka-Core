package controller

import (
	"fmt"
	"net/http"

	dto "Kavka/app/dto"
	"Kavka/app/presenters"
	"Kavka/internal/service"
	"Kavka/pkg/session"
	"Kavka/utils/bearer"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService}
}

func (ctrl *UserController) HandleLogin(ctx *gin.Context) {
	body := dto.Validate[dto.UserLoginDto](ctx)
	phone := body.Phone

	otp, err := ctrl.userService.Login(phone)
	if err != nil {
		presenters.ResponseError(ctx, err)
		return
	}

	// FIXME - Gonna be removed after implementing SMS service.
	fmt.Printf("OTP Code: %d\n", otp)

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
	refreshToken := ctx.GetHeader("refresh") //nolint

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
		userInfo, err := ctrl.userService.Authenticate(accessToken)
		if err != nil {
			presenters.ResponseError(ctx, err)
			return
		}

		presenters.ResponseUserInfo(ctx, userInfo)
	}
}
