package socket

import (
	"Kavka/internal/service"
	"Kavka/utils/bearer"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var clients []*websocket.Conn
var handlers = []func(MessageHandlerArgs) bool{
	NewChatsHandler,
	NewMessagesHandler,
}

type SocketService struct {
	userService *service.UserService
}

type SocketMessage struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

type MessageHandlerArgs struct {
	message       *SocketMessage
	conn          *websocket.Conn
	staticID      string
	socketService *SocketService
}

var upgrader = websocket.Upgrader{}

func NewSocketService(app *gin.Engine, userService *service.UserService) *SocketService {
	socketService := &SocketService{userService}

	app.GET("/ws", socketService.handleWebsocket)

	return socketService
}
func (s *SocketService) handleWebsocket(ctx *gin.Context) {
	// Authenticate
	accessToken, bearerOk := bearer.AccessToken(ctx)

	var staticID primitive.ObjectID

	if bearerOk {
		userInfo, err := s.userService.Authenticate(accessToken)
		if err != nil {
			ctx.Next()
			return
		}

		staticID = userInfo.StaticID
	}

	// Upgrade
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Next()
		return
	}

	defer conn.Close()

	for {
		var msgData *SocketMessage

		if err := conn.ReadJSON(&msgData); err != nil {
			log.Println("Unmarshal json error in socket:", err)
			break
		}

		clients = append(clients, conn)
		s.handleMessages(&MessageHandlerArgs{msgData, conn, staticID.String(), s})
	}
}

func (s *SocketService) handleMessages(args *MessageHandlerArgs) {
	var handled bool = false

	for _, handler := range handlers {
		result := handler(*args)
		if result {
			handled = true
			break
		}
	}

	if !handled {
		args.conn.WriteJSON(struct {
			Message string
		}{
			Message: "Invalid event",
		})
	}
}
