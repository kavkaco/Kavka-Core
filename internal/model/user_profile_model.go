package model

import "time"

type Profile struct {
	Biography     string         `bson:"biography" json:"biography"`
	ProfilePhotos []ProfilePhoto `bson:"profile_photos" json:"profilePhotos"`
	Devices       []Device       `bson:"devices" json:"devices"`
	BlockedUsers  []string       `bson:"blocked_users" json:"blockedUsers"`
	LastSeen      time.Time      `bson:"last_seen" json:"lastSeen"`
}

type ProfilePhoto struct {
	ID    string `bson:"id" json:"id"`
	Name  string `bson:"name" json:"name"`
	Order int    `bson:"order" json:"order"`
}

type Device struct {
	ID       string `bson:"id" json:"id"`
	Name     string `bson:"name" json:"name"`
	IP       string `bson:"ip" json:"ip"`
	Location string `bson:"location" json:"location"`
}
