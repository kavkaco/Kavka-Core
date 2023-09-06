package socket

func NewChatsHandler(args MessageHandlerArgs) bool {
	event := args.message.Event

	switch event {
	case "new_chat":
		return NewChat(args)
	}

	return false
}

func NewChat(args MessageHandlerArgs) bool {
	// username := args.message.Data["Username"]

	// Search in channels & groups
	// TODO

	return true
}
