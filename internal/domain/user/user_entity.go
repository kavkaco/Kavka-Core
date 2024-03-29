package user

import (
	"errors"
	"time"

	"github.com/kavkaco/Kavka-Core/pkg/session"
	"github.com/kavkaco/Kavka-Core/utils/random"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrEmptyPassword = errors.New("empty password")

type User struct {
	StaticID  primitive.ObjectID `bson:"id"        json:"id"`
	Name      string             `bson:"name" json:"name"`
	LastName  string             `bson:"last_name"  json:"lastName"`
	Phone     string             `bson:"phone" json:"phone"`
	Username  string             `bson:"username" json:"username"`
	Banned    bool               `bson:"banned" json:"banned"`
	Profile   Profile            `bson:"profile" json:"profile"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

func NewUser(phone string) *User {
	u := User{}
	u.Phone = phone
	u.StaticID = primitive.NewObjectID()
	u.Username = random.GenerateUsername()

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

// Interfaces

type UserRepository interface {
	Create(user *User) (*User, error)
	FindOne(filter bson.M) (*User, error)
	FindMany(filter bson.M) ([]*User, error)
	FindByID(staticID primitive.ObjectID) (*User, error)
	FindByUsername(username string) (*User, error)
	FindByPhone(phone string) (*User, error)
}

type Service interface {
	Login(phone string) error
	VerifyOTP(phone string, otp int) (*session.LoginTokens, error)
	RefreshToken(refreshToken string, accessToken string) (string, error)
	Authenticate(accessToken string) (*User, error)
}
