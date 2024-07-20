package stream

type MessagePayload struct {
	Data          interface{} `json:"data"`
	ReceiversUUID []string    `json:"receiversUUID"`
}
