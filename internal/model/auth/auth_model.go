package auth

import (
	"time"

	"github.com/kavkaco/Kavka-Core/internal/model/user"
)

type Auth struct {
	UserID              user.UserID    `bson:"user_id" json:"user_id"`
	PasswordHash        string         `bson:"password_hash"`
	FailedLoginAttempts int            `bson:"failed_login_attempts" json:"failed_login_attempts"`
	AccountLockedFor    *time.Duration `bson:"account_locked_for" json:"accountLockedFor"`
}
