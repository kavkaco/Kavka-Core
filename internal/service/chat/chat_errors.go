package chat

import "errors"

var (
	ErrInvalidValidation = errors.New("failed to validate arguments")

	ErrCreateChat                 = errors.New("unable to create chat")
	ErrNotFound                   = errors.New("not found")
	ErrUserNotFound               = errors.New("user not found")
	ErrGetUserChats               = errors.New("failed to get user chats")
	ErrChatAlreadyExists          = errors.New("chat already exists")
	ErrUnableToAddChatToUsersList = errors.New("unable to add add chat to users list")
	ErrUserJoinedBefore           = errors.New("user already is a member of the chat")
	ErrJoinDirectChat             = errors.New("joining direct chat is not possible")
	ErrRecipientNotFound          = errors.New("recipient not found")
	ErrMessageStoreCreation       = errors.New("failed to create message store for chat")
	ErrPublishEvent               = errors.New("failed to publish event to streamer")
)
