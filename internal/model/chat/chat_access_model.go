package chat

import (
	"slices"

	"github.com/kavkaco/Kavka-Core/internal/model/message"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// This function is used to check if a user with a specific `userStaticID` has access to send messages in a group chat.
// Being a member of a group is enough to have access to send the message for users.
func (detail *GroupChatDetail) HasAccessToSendMessage(userStaticID primitive.ObjectID) bool {
	return slices.Contains(detail.Members, userStaticID)
}

// This function is used to check if a user with a specific `userStaticID` has access to send messages in a channel chat.
// Only admins of the chat can send messages.
func (detail *ChannelChatDetail) HasAccessToSendMessage(userStaticID primitive.ObjectID) bool {
	return slices.Contains(detail.Admins, userStaticID)
}

// This function is used to check if a user with a specific `userStaticID` has access to delete a message in a group chat.
// The user is only allowed to delete messages his/her own messages.
// Admins can delete any messages.
func (detail *GroupChatDetail) HasAccessToDeleteMessage(userStaticID primitive.ObjectID, msg *message.Message) bool {
	// If is his/her own message
	if msg.SenderID == userStaticID {
		return true
	}

	// If is admin
	return slices.Contains(detail.Admins, userStaticID)
}

// This function is used to check if a user with a specific `userStaticID` has access to delete messages in a channel chat.
// Only admins of the chat can delete messages.
func (detail *ChannelChatDetail) HasAccessToDeleteMessage(userStaticID primitive.ObjectID) bool {
	return slices.Contains(detail.Admins, userStaticID)
}
