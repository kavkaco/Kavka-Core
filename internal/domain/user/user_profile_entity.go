package user

type UserProfile struct {
	Biography     string
	ProfilePhotos *[]UserProfilePhoto
	Devices       []Device
	Contacts      []string // reference to UserStaticID
	BlockedUsers  []string // reference to UserStaticID
	Setting
}

type UserProfilePhoto struct {
	ID     string
	Name   string
	IsMain bool
}

type Device struct {
	ID       string
	Name     string
	IP       string
	Location string
}

// TODO
