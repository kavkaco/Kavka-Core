package adapters

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kavkaco/Kavka-Core/app/presenters"
	"github.com/kavkaco/Kavka-Core/socket"
	"github.com/kavkaco/Kavka-Core/socket/handlers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var ErrCastConnInterface = errors.New("unable to cast connection interface")

// TODO - Write redis pub/sub for websocket connections.
var clients []*websocket.Conn

type socketAdapter struct {
	logger *zap.Logger
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Allow connections from any origin (not recommended for production)
}

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

func NewWebsocketAdapter(logger *zap.Logger) socket.SocketAdapter {
	return &socketAdapter{logger}
}

func (s *socketAdapter) Handle(ctx *gin.Context, handleConn func(conn interface{})) error {
	// Upgrade
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Next()
		return err
	}

	clients = append(clients, conn)

	handleConn(conn)

	return nil
}

func (s *socketAdapter) WriteMessage(conn interface{}, msg interface{}) error {
	if wsc, ok := conn.(*websocket.Conn); ok {
		return wsc.WriteJSON(msg)
	}

	return ErrCastConnInterface
}

func (s *socketAdapter) HandleMessages(conn interface{}, handleMessage func(msg socket.IncomingSocketMessage)) error {
	for {
		var msgData *socket.IncomingSocketMessage
		if wsc, ok := conn.(*websocket.Conn); ok {
			if err := wsc.ReadJSON(&msgData); err != nil {
				s.logger.Error("Unmarshal json error in websocket: " + err.Error())
				break
			}

			handleMessage(*msgData)
		} else {
			s.logger.Error("HandleMessages error: " + ErrCastConnInterface.Error())
			break
		}
	}

	return nil
}
