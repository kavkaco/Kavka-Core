package socket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"github.com/kavkaco/Kavka-Core/internal/domain/message"
	"github.com/kavkaco/Kavka-Core/internal/domain/user"
	"github.com/kavkaco/Kavka-Core/utils/bearer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var (
	clients  []*websocket.Conn
	handlers = []func(MessageHandlerArgs) bool{
		NewChatsHandler,
		NewMessagesHandler,
	}
)

type Service struct {
	logger      *zap.Logger
	userService user.Service
	chatService chat.Service
	msgService  message.Service
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

func NewSocketService(logger *zap.Logger, app *gin.Engine, userService user.Service, chatService chat.Service, messageService message.Service) *Service {
	socketService := &Service{logger, userService, chatService, messageService}

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
			s.logger.Error("Unmarshal json error in socket: " + err.Error())
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
