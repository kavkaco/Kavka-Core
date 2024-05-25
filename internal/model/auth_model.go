package model

type Auth struct {
	UserID              UserID `bson:"user_id"`
	PasswordHash        string `bson:"password_hash"`
	FailedLoginAttempts int    `bson:"failed_login_attempts"`
	AccountLockedUntil  int64  `bson:"account_locked_until"`
	EmailVerified       bool   `bson:"email_verified"`
}

func NewAuth(userID UserID, passwordHash string) *Auth {
	return &Auth{UserID: userID, PasswordHash: passwordHash, FailedLoginAttempts: 0, AccountLockedUntil: 0, EmailVerified: false}
}
