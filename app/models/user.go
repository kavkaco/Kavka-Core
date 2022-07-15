package models

import (
	"Nexus/app/database"
	"context"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	StaticID            int           `bson:"static_id"`
	Name                string        `bson:"name"`
	Email               string        `bson:"email"`
	Username            string        `bson:"username"`
	Bio                 string        `bson:"bio"`
	LastSeen            string        `bson:"last_seen"`
	Chats               []interface{} `bson:"chats"`
	ProfileImages       []string      `bson:"profile_images"`
	VerificCode         int           `bson:"verific_code"`
	VerificTryCount     int           `bson:"verific_try_count"`
	VerificCodeExpire   int64         `bson:"verific_code_expire"`
	VerificLimitDate    int64         `bson:"verific_limit_date"`
	FirstLoginCompleted bool          `bson:"first_login_completed"`
}

func MakeUserStaticID() int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	max := 99999999
	min := 11111111

	return min + r.Intn(max-min)
}

func FindUserByEmail(email string) *User {
	var user *User

	database.UsersCollection.FindOne(context.TODO(), bson.D{
		primitive.E{Key: "email", Value: email},
	}).Decode(&user)

	return user
}

func (user *User) IncreaseTryCount() {
	database.UsersCollection.FindOneAndUpdate(context.TODO(),
		bson.D{
			primitive.E{Key: "static_id", Value: user.StaticID},
		},
		bson.D{
			primitive.E{
				Key: "$inc",
				Value: bson.D{
					primitive.E{
						Key:   "verific_try_count",
						Value: 1,
					},
				},
			},
		},
	)
}
