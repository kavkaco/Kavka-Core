package repo_sql

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SQLSearchRepository struct {
	db *sql.DB
}

func NewSQLSearchRepository(db *sql.DB) *SQLSearchRepository {
	return &SQLSearchRepository{db: db}
}

func (r *SQLSearchRepository) Search(ctx context.Context, input string, limit int) (*model.SearchResultDTO, error) {
	result := &model.SearchResultDTO{
		Chats: []model.ChatDTO{},
		Users: []model.User{},
	}

	likePattern := "%" + input + "%"

	userRows, err := r.db.QueryContext(ctx,
		`SELECT user_id, name, last_name, email, username, biography, chats_list_ids, profile_photos
		 FROM users WHERE name ILIKE $1 OR last_name ILIKE $1 OR username ILIKE $1 OR email ILIKE $1
		 LIMIT $2`, likePattern, limit)
	if err == nil {
		defer userRows.Close()
		for userRows.Next() {
			user, err := r.scanUser(userRows)
			if err == nil {
				result.Users = append(result.Users, *user)
			}
		}
	}

	chatRows, err := r.db.QueryContext(ctx,
		`SELECT chat_id, chat_type, chat_detail FROM chats
		 WHERE chat_type IN ('group', 'channel')
		 AND (chat_detail::text ILIKE $1)
		 LIMIT $2`, likePattern, limit)
	if err == nil {
		defer chatRows.Close()
		for chatRows.Next() {
			var chatIDStr, chatType, detailJSON string
			err := chatRows.Scan(&chatIDStr, &chatType, &detailJSON)
			if err != nil {
				continue
			}

			chatID, _ := primitive.ObjectIDFromHex(chatIDStr)
			detail, _ := unmarshalChatDetail(chatType, detailJSON)

			result.Chats = append(result.Chats, model.ChatDTO{
				ChatID:     chatID,
				ChatType:   chatType,
				ChatDetail: detail,
			})
		}
	}

	return result, nil
}

func (r *SQLSearchRepository) SearchInChat(ctx context.Context, chatID model.ChatID, input string) ([]*model.MessageGetter, error) {
	likePattern := "%" + input + "%"

	rows, err := r.db.QueryContext(ctx,
		`SELECT message_id, chat_id, sender_id, sender_name, sender_last_name, sender_username,
		        created_at, edited, seen, type, content
		 FROM messages_v2 WHERE chat_id = $1 AND content::text ILIKE $2
		 ORDER BY created_at DESC LIMIT 50`,
		chatID.Hex(), likePattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	msgRepo := &SQLMessageRepository{db: r.db}
	var result []*model.MessageGetter
	for rows.Next() {
		doc, err := msgRepo.scanMessageV2(rows)
		if err != nil {
			continue
		}
		result = append(result, doc.ToMessageGetter())
	}

	return result, nil
}

func (r *SQLSearchRepository) scanUser(row interface{ Scan(...interface{}) error }) (*model.User, error) {
	var user model.User
	var chatsJSON, photosJSON string

	err := row.Scan(&user.UserID, &user.Name, &user.LastName, &user.Email,
		&user.Username, &user.Biography, &chatsJSON, &photosJSON)
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(chatsJSON), &user.ChatsListIDs)
	json.Unmarshal([]byte(photosJSON), &user.ProfilePhotos)
	if user.ChatsListIDs == nil {
		user.ChatsListIDs = []model.ChatID{}
	}

	return &user, nil
}
