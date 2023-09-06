package user

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// define errors
var (
	ErrEmptyPassword = errors.New("empty password")
)

type User struct {
	StaticID  primitive.ObjectID `bson:"_id" json:"static_id"`
	Name      string             `json:"name"`
	LastName  string             `bson:"last_name" json:"last_name"`
	Phone     string             `json:"phone"`
	Username  string             `json:"username"`
	Banned    bool               `json:"banned"`
	Profile   UserProfile        `json:"profile"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func NewUser(phone string) *User {
	u := User{}
	u.Phone = phone
	u.StaticID = primitive.NewObjectID()

	// set timestamps
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	return &u
}

func (u *User) FullName() string {
	return u.Name + " " + u.LastName
}

func (u User) IsBanned() bool {
	return u.Banned
}
