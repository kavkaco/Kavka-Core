package socket

import (
	"github.com/gin-gonic/gin"
)

type IncomingSocketMessage struct {
	Event string
	Data  map[string]interface{}
}

type OutgoingSocketMessage struct {
	Status int
	Event  string
	Data   interface{}
}

// TODO - rename to Communication Layer Adapter
type SocketAdapter interface {
	Handle(ctx *gin.Context, handleConn func(conn interface{})) error
	HandleMessages(conn interface{}, handleMessage func(msg IncomingSocketMessage)) error
	WriteMessage(conn interface{}, msg interface{}) error
}
