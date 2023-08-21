package user

import "time"

type UserProfile struct {
	Biography     string
	ProfilePhotos []UserProfilePhoto
	Devices       []Device
	Contacts      []string // reference to UserStaticID
	BlockedUsers  []string // reference to UserStaticID
	Setting
	LastSeen time.Time
}

type UserProfilePhoto struct {
	ID    string
	Name  string
	Order int
}

type Device struct {
	ID       string
	Name     string
	IP       string
	Location string
}
