package models

type DirectChatSides struct {
	User1StaticID int `bson:"user1_static_id"`
	User2StaticID int `bson:"user2_static_id"`
}

type DirectChat struct {
	ChatID   int             `bson:"chat_id"`
	Sides    DirectChatSides `bson:"chat_sides"`
	ChatType string          `bson:"chat_type"`
	Messages []Message       `bson:"messages"`
}

// Group or Channel
type Chat struct {
	ChatID          int       `bson:"chat_id"`
	ChatType        string    `bson:"chat_type"`
	Messages        []Message `bson:"messages"`
	ChatName        string    `bson:"chat_name"`
	ChatUsername    string    `bson:"chat_username"`
	ChatBio         string    `bson:"chat_bio"`
	ProfileImages   []string  `bson:"profile_images"`
	CreatorStaticID int       `bson:"creator_static_id"`
	Members         []int     `bson:"members"`
	Admins          []Admin   `bson:"members"`
}

type Admin struct {
	StaticID           int  `bson:"static_id"`
	CanChangeGroupInfo bool `bson:"can_change_group_info"`
	CanDeleteMessage   bool `bson:"can_delete_message"`
	CanBanUsers        bool `bson:"can_ban_users"`
	CanPinMessages     bool `bson:"can_pin_messages"`
	CanAddNewAdmins    bool `bson:"can_add_new_admins"`
}
