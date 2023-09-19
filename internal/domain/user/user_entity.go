package user

import (
	"errors"
	"time"

	"github.com/kavkaco/Kavka-Core/pkg/session"
	"github.com/kavkaco/Kavka-Core/utils/random"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// define errors.
var (
	ErrEmptyPassword = errors.New("empty password")
)

type User struct {
	StaticID  primitive.ObjectID `bson:"_id"        json:"static_id"`
	Name      string             `bson:"name" json:"name"`
	LastName  string             `bson:"last_name"  json:"last_name"`
	Phone     string             `bson:"phone" json:"phone"`
	Username  string             `bson:"username" json:"username"`
	Banned    bool               `bson:"banned" json:"banned"`
	Profile   Profile            `bson:"profile" json:"profile"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
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
	Where(filter bson.M) ([]*User, error)
	FindByID(staticID primitive.ObjectID) (*User, error)
	FindByUsername(username string) (*User, error)
	FindByPhone(phone string) (*User, error)
}

type UserService interface {
	Login(phone string) (int, error)
	VerifyOTP(phone string, otp int) (*session.LoginTokens, error)
	RefreshToken(refreshToken string, accessToken string) (string, error)
	Authenticate(accessToken string) (*User, error)
}
