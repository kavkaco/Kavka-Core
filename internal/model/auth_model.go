package model

import (
	"time"
)

type Auth struct {
	UserID              UserID         `bson:"user_id" json:"user_id"`
	PasswordHash        string         `bson:"password_hash"`
	FailedLoginAttempts int            `bson:"failed_login_attempts" json:"failed_login_attempts"`
	AccountLockedUntil  *time.Duration `bson:"account_locked_until" json:"accountLockedUntil"`
	EmailVerified       bool           `bson:"email_verified" json:"emailVerified"`
}

func NewAuth(userID UserID, passwordHash string) *Auth {
	return &Auth{UserID: userID, PasswordHash: passwordHash, FailedLoginAttempts: 0, AccountLockedUntil: nil, EmailVerified: false}
}
