package stream

// Set uuid to deliver the message to a specific user or set empty to broadcast to everyone
type Receiver struct {
	UUID string `json:"uuid"`
}

type MessagePayload struct {
	Data      map[string]interface{} `json:"data"`
	Receivers []Receiver             `json:"receivers"`
}
