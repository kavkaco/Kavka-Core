package model

import (
	"slices"
)

// This function is used to check if a user with a specific `userID` has access to send messages in a group chat.
// Being a member of a group is enough to have access to send the message for users.
func (detail *GroupChatDetail) HasAccessToSendMessage(userID UserID) bool {
	return slices.Contains(detail.Members, userID)
}

// This function is used to check if a user with a specific `userID` has access to send messages in a channel chat.
// Only admins of the chat can send messages.
func (detail *ChannelChatDetail) HasAccessToSendMessage(userID UserID) bool {
	return slices.Contains(detail.Admins, userID)
}

// This function is used to check if a user with a specific `userID` has access to delete a message in a group chat.
// The user is only allowed to delete messages his/her own messages.
// Admins can delete any messages.
func (detail *GroupChatDetail) HasAccessToDeleteMessage(userID UserID, msg *Message) bool {
	// If is his/her own message
	if msg.SenderID == userID {
		return true
	}

	// If is admin
	return slices.Contains(detail.Admins, userID)
}

// This function is used to check if a user with a specific `userID` has access to delete messages in a channel chat.
// Only admins of the chat can delete messages.
func (detail *ChannelChatDetail) HasAccessToDeleteMessage(userID UserID) bool {
	return slices.Contains(detail.Admins, userID)
}
