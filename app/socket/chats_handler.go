package socket

import "log"

func NewChatsHandler(args MessageHandlerArgs) bool {
	event := args.message.Event

	switch event {
	case "new_chat":
		return NewChat(args)
	}

	return false
}

func NewChat(args MessageHandlerArgs) bool {
	username := args.message.Data["username"]

	log.Println(username)

	// Search in channels & groups
	// TODO

	return true
}
