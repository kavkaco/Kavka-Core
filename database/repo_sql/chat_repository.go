package repo_sql

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SQLChatRepository struct {
	db *sql.DB
}

func NewSQLChatRepository(db *sql.DB) *SQLChatRepository {
	return &SQLChatRepository{db: db}
}

func (r *SQLChatRepository) Create(ctx context.Context, chat model.Chat) (*model.Chat, error) {
	detailJSON, err := json.Marshal(chat.ChatDetail)
	if err != nil {
		return nil, err
	}

	_, err = r.db.ExecContext(ctx,
		`INSERT INTO chats (chat_id, chat_type, chat_detail) VALUES ($1, $2, $3)`,
		chat.ChatID.Hex(), chat.ChatType, string(detailJSON))
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *SQLChatRepository) GetChat(ctx context.Context, chatID model.ChatID) (*model.Chat, error) {
	var chat model.Chat
	var chatIDStr, detailJSON string

	err := r.db.QueryRowContext(ctx,
		`SELECT chat_id, chat_type, chat_detail FROM chats WHERE chat_id = $1`, chatID.Hex()).
		Scan(&chatIDStr, &chat.ChatType, &detailJSON)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	chat.ChatID, _ = primitive.ObjectIDFromHex(chatIDStr)
	chat.ChatDetail, _ = unmarshalChatDetail(chat.ChatType, detailJSON)

	return &chat, nil
}

func (r *SQLChatRepository) Destroy(ctx context.Context, chatID model.ChatID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM chats WHERE chat_id = $1`, chatID.Hex())
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return repository.ErrNotModified
	}
	return nil
}

func (r *SQLChatRepository) GetUserChats(ctx context.Context, userID model.UserID, chatIDs []model.ChatID) ([]model.ChatDTO, error) {
	if len(chatIDs) == 0 {
		return []model.ChatDTO{}, nil
	}

	hexIDs := make([]string, 0, len(chatIDs))
	idMap := make(map[string]model.ChatID, len(chatIDs))
	for _, id := range chatIDs {
		h := id.Hex()
		hexIDs = append(hexIDs, h)
		idMap[h] = id
	}

	query := `SELECT chat_id, chat_type, chat_detail FROM chats WHERE chat_id IN (`
	args := make([]interface{}, len(hexIDs))
	for i, h := range hexIDs {
		if i > 0 {
			query += ","
		}
		query += `$` + string(rune('1'+i))
		args[i] = h
	}
	query += `)`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []model.ChatDTO
	for rows.Next() {
		var chatIDStr, chatType, detailJSON string
		err := rows.Scan(&chatIDStr, &chatType, &detailJSON)
		if err != nil {
			return nil, err
		}

		detail, _ := unmarshalChatDetail(chatType, detailJSON)
		chatID := idMap[chatIDStr]

		lastMessage, _ := r.fetchLastMessageForChat(ctx, chatID)

		chats = append(chats, model.ChatDTO{
			ChatID:      chatID,
			ChatType:    chatType,
			ChatDetail:  detail,
			LastMessage: lastMessage,
		})
	}

	return chats, nil
}

func (r *SQLChatRepository) GetDirectChat(ctx context.Context, userID, recipientUserID model.UserID) (*model.Chat, error) {
	var chat model.Chat
	var chatIDStr, chatType, detailJSON string

	err := r.db.QueryRowContext(ctx,
		`SELECT chat_id, chat_type, chat_detail FROM chats
		 WHERE chat_type = 'direct'
		 AND (chat_detail::text LIKE $1 OR chat_detail::text LIKE $2)
		 LIMIT 1`,
		`%"user_id":"`+userID+`"%"recipient_user_id":"`+recipientUserID+`"`,
		`%"user_id":"`+recipientUserID+`"%"recipient_user_id":"`+userID+`"`,
	).Scan(&chatIDStr, &chatType, &detailJSON)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	chat.ChatID, _ = primitive.ObjectIDFromHex(chatIDStr)
	chat.ChatType = chatType
	chat.ChatDetail, _ = unmarshalChatDetail(chatType, detailJSON)

	return &chat, nil
}

func (r *SQLChatRepository) GetChatMembers(ctx context.Context, chatID model.ChatID) ([]model.Member, error) {
	chat, err := r.GetChat(ctx, chatID)
	if err != nil {
		return nil, err
	}

	switch d := chat.ChatDetail.(type) {
	case *model.ChannelChatDetail:
		return r.membersFromIDs(ctx, d.Members)
	case *model.GroupChatDetail:
		return r.membersFromIDs(ctx, d.Members)
	default:
		return []model.Member{}, nil
	}
}

func (r *SQLChatRepository) membersFromIDs(ctx context.Context, userIDs []model.UserID) ([]model.Member, error) {
	if len(userIDs) == 0 {
		return []model.Member{}, nil
	}

	query := `SELECT user_id, name, last_name FROM users WHERE user_id IN (`
	args := make([]interface{}, len(userIDs))
	for i, uid := range userIDs {
		if i > 0 {
			query += ","
		}
		query += `$` + string(rune('1'+i))
		args[i] = uid
	}
	query += `)`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []model.Member
	for rows.Next() {
		var m model.Member
		if err := rows.Scan(&m.UserID, &m.Name, &m.LastName); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *SQLChatRepository) JoinChat(ctx context.Context, chatType string, userID string, chatID model.ChatID) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_chats (user_id, chat_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, chatID.Hex())
	return err
}

func (r *SQLChatRepository) AddToUsersChatsList(ctx context.Context, userID string, chatID model.ChatID) error {
	return r.JoinChat(ctx, "", userID, chatID)
}

func (r *SQLChatRepository) fetchLastMessageForChat(ctx context.Context, chatID model.ChatID) (*model.Message, error) {
	var msg model.Message
	var msgIDStr, senderID, msgType, contentJSON string

	err := r.db.QueryRowContext(ctx,
		`SELECT message_id, sender_id, created_at, edited, seen, type, content
		 FROM messages_v2 WHERE chat_id = $1 ORDER BY created_at DESC LIMIT 1`, chatID.Hex()).
		Scan(&msgIDStr, &senderID, &msg.CreatedAt, &msg.Edited, &msg.Seen, &msgType, &contentJSON)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	msg.MessageID, _ = primitive.ObjectIDFromHex(msgIDStr)
	msg.SenderID = senderID
	msg.Type = msgType
	json.Unmarshal([]byte(contentJSON), &msg.Content)

	return &msg, nil
}

func unmarshalChatDetail(chatType string, jsonData string) (interface{}, error) {
	switch chatType {
	case "channel":
		var d model.ChannelChatDetail
		err := json.Unmarshal([]byte(jsonData), &d)
		return &d, err
	case "group":
		var d model.GroupChatDetail
		err := json.Unmarshal([]byte(jsonData), &d)
		return &d, err
	case "direct":
		var d model.DirectChatDetail
		err := json.Unmarshal([]byte(jsonData), &d)
		return &d, err
	}
	return nil, nil
}
