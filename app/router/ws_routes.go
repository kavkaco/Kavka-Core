package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/kavkaco/Kavka-Core/app/presenters"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/socket"
	"github.com/kavkaco/Kavka-Core/socket/handlers"
	"go.uber.org/zap"
)

func WebsocketRoute(ctx context.Context, logger *zap.Logger, websocketAdapter socket.SocketAdapter, handlerServices handlers.HandlerServices) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// Get UserStaticID form AuthenticatedMiddleware and cast to primitive.ObjectId!
		userIDAny, getOk := ctx.Get("user_static_id")

		if getOk {
			userID, castOk := userIDAny.(model.UserID)
			if !castOk {
				logger.Error("Unable to cast any to string")
				ctx.Next()
				return
			}

			// Call handle from WebsocketAdapter and pass the conn to the handler
			err := websocketAdapter.Handle(ctx, func(conn interface{}) {
				handlerErr := handlers.NewSocketHandler(ctx, logger, websocketAdapter, conn, &handlerServices, userID)
				if handlerErr != nil {
					logger.Error("Unable to create socket handler instance")
				}
			})
			if err != nil {
				presenters.InternalServerErrorResponse(ctx)
			}
		} else {
			logger.Error("Unable to read user_static_id from gin.Context")
			ctx.Next()
		}
	}
}
