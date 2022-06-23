package models

import "time"

type User struct {
	StaticID             uint          `bson:"static_id"`
	Name                 string        `bson:"name"`
	Email                string        `bson:"email"`
	Username             string        `bson:"username"`
	Bio                  string        `bson:"bio"`
	LastSeen             string        `bson:"last_seen"`
	Chats                []interface{} `bson:"chats"`
	ProfilePhotos        []string      `bson:"profile_photos"`
	VerificCode          uint          `bson:"verific_code"`
	VerificTryCount      uint          `bson:"verific_try_count"`
	VerificCodeExpire    time.Time     `bson:"verific_code_expire"`
	VerificCodeLimitDate time.Time     `bson:"verific_code_limit_date"`
	FirstLoginCompleted  bool          `bson:"first_login_completed"`
}
