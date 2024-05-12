package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kavkaco/Kavka-Core/app/presenters"
	"github.com/kavkaco/Kavka-Core/socket"
	"github.com/kavkaco/Kavka-Core/socket/handlers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func WebsocketRoute(logger *zap.Logger, websocketAdapter socket.SocketAdapter, handlerServices handlers.HandlerServices) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// Get UserStaticID form AuthenticatedMiddleware and cast to primitive.ObjectId!
		if userStaticIDAny, ok := ctx.Get("user_static_id"); ok {
			userStaticIDStr, _ := userStaticIDAny.(string)

			userStaticID, err := primitive.ObjectIDFromHex(userStaticIDStr)
			if err != nil {
				logger.Error("Unable to cast string to ObjectId")
				ctx.Next()
			}

			// Call handle from WebsocketAdapter and pass the conn to the handler
			err = websocketAdapter.Handle(ctx, func(conn interface{}) {
				handlerErr := handlers.NewSocketHandler(logger, websocketAdapter, conn, &handlerServices, userStaticID)
				if handlerErr != nil {
					logger.Error("Unable to create SocketHandler instance: " + err.Error())
				}
			})
			if err != nil {
				presenters.ResponseInternalServerError(ctx)
			}
		} else {
			logger.Error("Unable to read user_static_id from gin.Context")
			ctx.Next()
		}
	}
}
