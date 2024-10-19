package model

type AccountDetail struct {
	Devices      []Device `bson:"devices" json:"devices"`
	BlockedUsers []string `bson:"blocked_users" json:"blockedUsers"`
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

type Setting struct {
	ShowEmail    bool `bson:"show_email" json:"showEmail"`
	ShowLastSeen bool `bson:"show_last_seen" json:"showLastSeen"`
	TwoFA        bool `bson:"two_fa" json:"twoFa"`
}
