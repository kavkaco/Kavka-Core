package repository

import (
	"context"
	"time"

	"github.com/kavkaco/Kavka-Core/internal/model"
)

type AuthRepository interface {
	Create(ctx context.Context, authModel *model.Auth) (*model.Auth, error)
	GetUserAuth(ctx context.Context, userID model.UserID) (*model.Auth, error)
	ChangePassword(ctx context.Context, userID model.UserID, passwordHash string) error
	VerifyEmail(ctx context.Context, userID model.UserID) error
	IncrementFailedLoginAttempts(ctx context.Context, userID model.UserID) error
	ClearFailedLoginAttempts(ctx context.Context, userID model.UserID) error
	LockAccount(ctx context.Context, userID model.UserID, lockDuration time.Duration) error
	UnlockAccount(ctx context.Context, userID model.UserID) error
}
