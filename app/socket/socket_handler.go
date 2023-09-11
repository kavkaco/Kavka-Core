package socket

import (
	"log"

	"Kavka/internal/service"
	"Kavka/utils/bearer"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	clients  []*websocket.Conn
	handlers = []func(MessageHandlerArgs) bool{
		NewChatsHandler,
		NewMessagesHandler,
	}
)

type Service struct {
	userService *service.UserService
	chatService *service.ChatService
	msgService  *service.MessageService
}

type Message struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

type MessageHandlerArgs struct {
	message       *Message
	conn          *websocket.Conn
	staticID      primitive.ObjectID
	socketService *Service
}

var upgrader = websocket.Upgrader{}

func NewSocketService(app *gin.Engine, userService *service.UserService,
	chatService *service.ChatService, msgService *service.MessageService,
) *Service {
	socketService := &Service{userService, chatService, msgService}

	app.GET("/ws", socketService.handleWebsocket)

	return socketService
}

func (s *Service) handleWebsocket(ctx *gin.Context) {
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
		var msgData *Message

		if err := conn.ReadJSON(&msgData); err != nil {
			log.Println("Unmarshal json error in socket:", err)
			break
		}

		clients = append(clients, conn)
		s.handleMessages(&MessageHandlerArgs{msgData, conn, staticID, s})
	}
}

func (s *Service) handleMessages(args *MessageHandlerArgs) {
	handled := false

	for _, handler := range handlers {
		result := handler(*args)
		if result {
			handled = true
			break
		}
	}

	if !handled {
		args.conn.WriteJSON(struct { //nolint
			Message string
		}{
			Message: "Invalid event",
		})
	}
}
