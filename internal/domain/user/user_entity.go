package user

import (
	"errors"
	"time"

	"Kavka/pkg/uuid"
)

// define errors
var (
	ErrEmptyPassword = errors.New("empty password")
)

type User struct {
	StaticID  string
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
	u.StaticID = u.NewStaticID()
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

func (u User) NewStaticID() string {
	return uuid.Random()
}
