package repo_sql

import (
	"context"
	"database/sql"
	"time"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
)

type SQLAuthRepository struct {
	db *sql.DB
}

func NewSQLAuthRepository(db *sql.DB) *SQLAuthRepository {
	return &SQLAuthRepository{db: db}
}

func (r *SQLAuthRepository) Create(ctx context.Context, auth *model.Auth) (*model.Auth, error) {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_auth (user_id, password_hash, failed_login_attempts, account_locked_until, email_verified)
		 VALUES ($1, $2, $3, $4, $5)`,
		auth.UserID, auth.PasswordHash, auth.FailedLoginAttempts, auth.AccountLockedUntil, boolToInt(auth.EmailVerified))
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (r *SQLAuthRepository) GetUserAuth(ctx context.Context, userID model.UserID) (*model.Auth, error) {
	var auth model.Auth
	var emailVerified int

	err := r.db.QueryRowContext(ctx,
		`SELECT user_id, password_hash, failed_login_attempts, account_locked_until, email_verified
		 FROM user_auth WHERE user_id = $1`, userID).
		Scan(&auth.UserID, &auth.PasswordHash, &auth.FailedLoginAttempts, &auth.AccountLockedUntil, &emailVerified)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	auth.EmailVerified = emailVerified != 0
	return &auth, nil
}

func (r *SQLAuthRepository) ChangePassword(ctx context.Context, userID model.UserID, passwordHash string) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE user_auth SET password_hash = $1 WHERE user_id = $2`, passwordHash, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *SQLAuthRepository) VerifyEmail(ctx context.Context, userID model.UserID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE user_auth SET email_verified = 1 WHERE user_id = $1`, userID)
	return err
}

func (r *SQLAuthRepository) IncrementFailedLoginAttempts(ctx context.Context, userID model.UserID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE user_auth SET failed_login_attempts = failed_login_attempts + 1 WHERE user_id = $1`, userID)
	return err
}

func (r *SQLAuthRepository) ClearFailedLoginAttempts(ctx context.Context, userID model.UserID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE user_auth SET failed_login_attempts = 0 WHERE user_id = $1`, userID)
	return err
}

func (r *SQLAuthRepository) LockAccount(ctx context.Context, userID model.UserID, lockDuration time.Duration) error {
	lockedUntil := time.Now().Add(lockDuration).Unix()
	_, err := r.db.ExecContext(ctx,
		`UPDATE user_auth SET account_locked_until = $1 WHERE user_id = $2`, lockedUntil, userID)
	return err
}

func (r *SQLAuthRepository) UnlockAccount(ctx context.Context, userID model.UserID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE user_auth SET account_locked_until = 0 WHERE user_id = $1`, userID)
	return err
}

func (r *SQLAuthRepository) DeleteByID(ctx context.Context, userID model.UserID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM user_auth WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return repository.ErrNotDeleted
	}
	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
