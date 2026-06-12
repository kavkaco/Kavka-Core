package repo_sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
)

type SQLUserRepository struct {
	db *sql.DB
}

func NewSQLUserRepository(db *sql.DB) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	chatsJSON, _ := json.Marshal(user.ChatsListIDs)
	photosJSON, _ := json.Marshal(user.ProfilePhotos)

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (user_id, name, last_name, email, username, biography, chats_list_ids, profile_photos)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.UserID, user.Name, user.LastName, user.Email, user.Username, user.Biography,
		string(chatsJSON), string(photosJSON))
	if err != nil {
		if isUniqueViolation(err) {
			return nil, repository.ErrUniqueConstraint
		}
		return nil, err
	}
	return user, nil
}

func (r *SQLUserRepository) FindByUserID(ctx context.Context, userID model.UserID) (*model.User, error) {
	return r.scanUser(ctx, `SELECT user_id, name, last_name, email, username, biography, chats_list_ids, profile_photos FROM users WHERE user_id = $1`, userID)
}

func (r *SQLUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return r.scanUser(ctx, `SELECT user_id, name, last_name, email, username, biography, chats_list_ids, profile_photos FROM users WHERE username = $1`, username)
}

func (r *SQLUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.scanUser(ctx, `SELECT user_id, name, last_name, email, username, biography, chats_list_ids, profile_photos FROM users WHERE email = $1`, email)
}

func (r *SQLUserRepository) scanUser(ctx context.Context, query string, args ...interface{}) (*model.User, error) {
	var user model.User
	var chatsJSON, photosJSON string

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.UserID, &user.Name, &user.LastName, &user.Email,
		&user.Username, &user.Biography, &chatsJSON, &photosJSON)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(chatsJSON), &user.ChatsListIDs)
	json.Unmarshal([]byte(photosJSON), &user.ProfilePhotos)
	if user.ChatsListIDs == nil {
		user.ChatsListIDs = []model.ChatID{}
	}

	return &user, nil
}

func (r *SQLUserRepository) Update(ctx context.Context, userID string, name, lastName, username, biography string) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE users SET name=$1, last_name=$2, username=$3, biography=$4 WHERE user_id=$5`,
		name, lastName, username, biography, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return repository.ErrNotModified
	}
	return nil
}

func (r *SQLUserRepository) DeleteByID(ctx context.Context, userID model.UserID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return repository.ErrNotDeleted
	}
	return nil
}

func (r *SQLUserRepository) GetChats(ctx context.Context, userID model.UserID) ([]model.ChatID, error) {
	user, err := r.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user.ChatsListIDs, nil
}

func (r *SQLUserRepository) AddToUserChats(ctx context.Context, userID model.UserID, chatID model.ChatID) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_chats (user_id, chat_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, chatID.Hex())
	return err
}

func (r *SQLUserRepository) IsIndexesUnique(ctx context.Context, email, username string) (bool, []string) {
	var conflicts []string

	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE email = $1 OR username = $2`, email, username).Scan(&count)
	if err != nil || count == 0 {
		return true, nil
	}

	var foundEmail, foundUsername string
	r.db.QueryRowContext(ctx, `SELECT email, username FROM users WHERE email = $1 OR username = $2 LIMIT 1`, email, username).Scan(&foundEmail, &foundUsername)

	if strings.EqualFold(foundEmail, email) {
		conflicts = append(conflicts, "email")
	}
	if strings.EqualFold(foundUsername, username) {
		conflicts = append(conflicts, "username")
	}

	return len(conflicts) == 0, conflicts
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE") || strings.Contains(msg, "unique") || strings.Contains(msg, "23505")
}
