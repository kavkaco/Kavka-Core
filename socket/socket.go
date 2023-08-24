package socket

import (
	"Kavka/service"
	"Kavka/utils/bearer"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SocketService struct {
	userService *service.UserService
}

type SocketMessage struct {
	Event string                 `json:"Event"`
	Data  map[string]interface{} `json:"Data"`
}

func (s *SocketService) handleMessages(message *SocketMessage, conn *websocket.Conn, staticID string) {
	NewMessagesHandler(message, conn.Conn, staticID)
}

func NewSocketService(app *fiber.App, userService *service.UserService) *SocketService {
	socketService := &SocketService{userService}

	app.Use("/ws", socketService.endpoint)
	app.Get("/ws", websocket.New(socketService.handleWebsocket))

	return socketService
}

func (s *SocketService) handleWebsocket(ctx *websocket.Conn) {
	staticID := ctx.Locals("StaticID").(primitive.ObjectID).Hex()

	for {
		var msgData *SocketMessage

		if err := ctx.ReadJSON(&msgData); err != nil {
			log.Println(err)
			break
		}

		s.handleMessages(msgData, ctx, staticID)
	}
}

func (s *SocketService) endpoint(ctx *fiber.Ctx) error {
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
