package model

type MessageSenderDTO struct {
	UserID   UserID `bson:"user_id" json:"userID"`
	Name     string `bson:"name" json:"name"`
	LastName string `bson:"last_name" json:"lastName"`
	Username string `bson:"username" json:"username"`
}
