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
	StaticID  primitive.ObjectID `bson:"_id"`
	Name      string
	LastName  string
	Phone     string
	Banned    bool
	Profile   UserProfile
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(phone string) *User {
	u := User{}
	u.Phone = phone
	u.StaticID = primitive.NewObjectID()
	u.Banned = false

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

type CreateUserData struct {
	Name     string
	LastName string
	Phone    string
}
