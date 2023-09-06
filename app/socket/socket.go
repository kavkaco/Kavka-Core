package socket

import (
	"Kavka/internal/service"
	"Kavka/utils/bearer"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var clients []*websocket.Conn

type SocketService struct {
	userService *service.UserService
}

type SocketMessage struct {
	Event string                 `json:"Event"`
	Data  map[string]interface{} `json:"Data"`
}

type MessageHandlerArgs struct {
	message       *SocketMessage
	conn          *websocket.Conn
	staticID      string
	socketService *SocketService
}

func NewSocketService(app *fiber.App, userService *service.UserService) *SocketService {
	socketService := &SocketService{userService}

	app.Use("/ws", socketService.endpoint)
	app.Get("/ws", websocket.New(socketService.handleWebsocket))

	return socketService
}

func (s *SocketService) handleMessages(args MessageHandlerArgs) {
	var handled bool = false

	handlers := []func(MessageHandlerArgs) bool{
		NewChatsHandler,
		NewMessagesHandler,
	}

	for _, handler := range handlers {
		result := handler(args)
		if result {
			handled = true
			return
		}
	}

	if !handled {
		args.conn.WriteJSON(struct {
			Message string
		}{
			Message: "Invalid Event",
		})
	}
}

func (s *SocketService) handleWebsocket(ctx *websocket.Conn) {
	staticID := ctx.Locals("StaticID").(primitive.ObjectID).Hex()

	for {
		var msgData *SocketMessage

		if err := ctx.ReadJSON(&msgData); err != nil {
			log.Println(err)
			break
		}

		clients = append(clients, ctx)
		s.handleMessages(MessageHandlerArgs{msgData, ctx, staticID, s})
	}
}

func (s *SocketService) endpoint(ctx *gin.Context) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		accessToken, bearerOk := bearer.AccessToken(ctx)

		if bearerOk {
			userInfo, err := s.userService.Authenticate(accessToken)
			if err != nil {
				return fiber.ErrUpgradeRequired
			}

			staticID := userInfo.StaticID
			phone := userInfo.Phone

			ctx.Locals("StaticID", staticID)
			ctx.Locals("Phone", phone)
			ctx.Locals("allowed", true)

			return ctx.Next()
		}
	}

	return fiber.ErrUpgradeRequired
}
