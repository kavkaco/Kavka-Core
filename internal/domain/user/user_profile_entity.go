package user

import "github.com/google/uuid"

type UserProfile struct {
	Biography     string
	ProfilePhotos *[]UserProfilePhoto
	Devices       []Device
	Contacts      []uuid.UUID
	BlockedUsers  []uuid.UUID
	Setting
}

type UserProfilePhoto struct {
	ID     uuid.UUID
	Name   string
	IsMain bool
}

type Device struct {
	ID       uuid.UUID
	Name     string
	IP       string
	Location string
}

// TODO
