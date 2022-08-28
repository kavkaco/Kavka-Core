package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// define errors
var (
	ErrEmptyPassword = errors.New("empty password")
)

type User struct {
	StaticID      uuid.UUID
	Name          string
	LastName      string
	Username      string
	PasswordHash  string
	Email         string
	EmailVerified bool
	Banned        bool
	Profile       UserProfile
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (u *User) FullName() string {
	return u.Name + " " + u.LastName
}

func (u User) IsVerified() bool {
	return u.EmailVerified
}

func (u User) IsBanned() bool {
	return u.Banned
}

func (u User) NewUserID() uuid.UUID {
	return uuid.New()
}

func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return ErrEmptyPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)

	return nil
}

func (u User) IsValidPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

func (u *User) PrepareToCreate() {
	u.StaticID = u.NewUserID()
	u.EmailVerified = false
	u.Banned = false
	// set timestamps
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
}
