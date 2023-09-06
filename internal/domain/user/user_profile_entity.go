package user

import "time"

type UserProfile struct {
	Biography     string             `json:"biography"`
	ProfilePhotos []UserProfilePhoto `json:"profile_photos"`
	Devices       []Device           `json:"devices"`
	Contacts      []string           `json:"contacts"`
	BlockedUsers  []string           `json:"blocked_users"`
	Setting       Setting            `json:"setting"`
	LastSeen      time.Time          `json:"last_seen"`
}

type UserProfilePhoto struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Order int    `json:"order"`
}

type Device struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IP       string `json:"ip"`
	Location string `json:"location"`
}
