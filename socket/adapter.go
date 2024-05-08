package socket

import (
	"github.com/gin-gonic/gin"
)

type SocketMessage struct {
	Event string
	Data  map[string]interface{}
}

type SocketAdapter interface {
	// OpenConnection establishes a connection to the socket server
	OpenConnection(app *gin.Engine, handleConn func(conn interface{})) error
	// CloseConnection closes the connection to the socket server
	CloseConnection() error
	// ReadMessage reads a message from the socket
	HandleMessages(conn interface{}, handleMessage func(msg SocketMessage)) error
	// WriteMessage writes a message to the socket
	WriteMessage(conn interface{}, msg interface{}) error
}
