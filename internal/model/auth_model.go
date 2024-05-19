package model

import (
	"time"
)

type Auth struct {
	UserID              UserID         `bson:"user_id" json:"user_id"`
	PasswordHash        string         `bson:"password_hash"`
	FailedLoginAttempts int            `bson:"failed_login_attempts" json:"failed_login_attempts"`
	AccountLockedFor    *time.Duration `bson:"account_locked_for" json:"accountLockedFor"`
}
