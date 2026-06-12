package repo_sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SQLMessageRepository struct {
	db *sql.DB
}

func NewSQLMessageRepository(db *sql.DB) *SQLMessageRepository {
	return &SQLMessageRepository{db: db}
}

func (r *SQLMessageRepository) Create(ctx context.Context, chatID model.ChatID) error {
	return nil
}

func (r *SQLMessageRepository) Insert(ctx context.Context, chatID model.ChatID, message *model.Message) (*model.Message, error) {
	contentJSON, _ := json.Marshal(message.Content)

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO messages_v2 (message_id, chat_id, sender_id, created_at, edited, seen, type, content)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		message.MessageID.Hex(), chatID.Hex(), message.SenderID,
		message.CreatedAt, boolToInt(message.Edited), boolToInt(message.Seen),
		message.Type, string(contentJSON))
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (r *SQLMessageRepository) FetchLastMessage(ctx context.Context, chatID model.ChatID) (*model.Message, error) {
	return r.fetchLastMessageForChat(ctx, chatID)
}

func (r *SQLMessageRepository) FetchMessage(ctx context.Context, chatID model.ChatID, messageID model.MessageID) (*model.Message, error) {
	doc, err := r.FetchMessageV2(ctx, chatID, messageID)
	if err != nil {
		return nil, err
	}
	return &model.Message{
		MessageID: doc.ID,
		SenderID:  doc.SenderID,
		CreatedAt: doc.CreatedAt,
		Edited:    doc.Edited,
		Seen:      doc.Seen,
		Type:      doc.Type,
		Content:   doc.Content,
	}, nil
}

func (r *SQLMessageRepository) FetchMessages(ctx context.Context, chatID model.ChatID) ([]*model.MessageGetter, error) {
	return r.FetchMessagesPaginated(ctx, chatID, 0, 200)
}

func (r *SQLMessageRepository) FetchMessagesPaginated(ctx context.Context, chatID model.ChatID, skip, limit int64) ([]*model.MessageGetter, error) {
	docs, err := r.FetchMessagesV2(ctx, chatID, skip, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*model.MessageGetter, 0, len(docs))
	for _, doc := range docs {
		result = append(result, doc.ToMessageGetter())
	}
	return result, nil
}

func (r *SQLMessageRepository) UpdateMessageContent(ctx context.Context, chatID model.ChatID, messageID model.MessageID, newMessageContent string) error {
	return r.UpdateMessageContentV2(ctx, chatID, messageID, newMessageContent)
}

func (r *SQLMessageRepository) Delete(ctx context.Context, chatID model.ChatID, messageID model.MessageID) error {
	return r.DeleteV2(ctx, chatID, messageID)
}

func (r *SQLMessageRepository) InsertV2(ctx context.Context, doc *model.MessageDocumentV2) (*model.MessageDocumentV2, error) {
	if doc.ID.IsZero() {
		doc.ID = primitive.NewObjectID()
	}
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = time.Now()
	}

	contentJSON, _ := json.Marshal(doc.Content)

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO messages_v2 (message_id, chat_id, sender_id, sender_name, sender_last_name, sender_username, created_at, edited, seen, type, content)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		doc.ID.Hex(), doc.ChatID.Hex(), doc.SenderID, doc.SenderName, doc.SenderLastName, doc.SenderUsername,
		doc.CreatedAt, boolToInt(doc.Edited), boolToInt(doc.Seen), doc.Type, string(contentJSON))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (r *SQLMessageRepository) FetchMessagesV2(ctx context.Context, chatID model.ChatID, skip, limit int64) ([]*model.MessageDocumentV2, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT message_id, chat_id, sender_id, sender_name, sender_last_name, sender_username,
		        created_at, edited, seen, type, content
		 FROM messages_v2 WHERE chat_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		chatID.Hex(), limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []*model.MessageDocumentV2
	for rows.Next() {
		doc, err := r.scanMessageV2(rows)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

func (r *SQLMessageRepository) FetchMessageV2(ctx context.Context, chatID model.ChatID, messageID model.MessageID) (*model.MessageDocumentV2, error) {
	doc, err := r.scanMessageV2(r.db.QueryRowContext(ctx,
		`SELECT message_id, chat_id, sender_id, sender_name, sender_last_name, sender_username,
		        created_at, edited, seen, type, content
		 FROM messages_v2 WHERE chat_id = $1 AND message_id = $2`,
		chatID.Hex(), messageID.Hex()))
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return doc, nil
}

func (r *SQLMessageRepository) FetchLastMessageV2(ctx context.Context, chatID model.ChatID) (*model.MessageDocumentV2, error) {
	doc, err := r.scanMessageV2(r.db.QueryRowContext(ctx,
		`SELECT message_id, chat_id, sender_id, sender_name, sender_last_name, sender_username,
		        created_at, edited, seen, type, content
		 FROM messages_v2 WHERE chat_id = $1 ORDER BY created_at DESC LIMIT 1`,
		chatID.Hex()))
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return doc, nil
}

func (r *SQLMessageRepository) UpdateMessageContentV2(ctx context.Context, chatID model.ChatID, messageID model.MessageID, newMessageContent string) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE messages_v2 SET content = json_set(content, '$.text', $1), edited = 1
		 WHERE chat_id = $2 AND message_id = $3`,
		newMessageContent, chatID.Hex(), messageID.Hex())
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *SQLMessageRepository) DeleteV2(ctx context.Context, chatID model.ChatID, messageID model.MessageID) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM messages_v2 WHERE chat_id = $1 AND message_id = $2`,
		chatID.Hex(), messageID.Hex())
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return repository.ErrNotModified
	}
	return nil
}

func (r *SQLMessageRepository) CountMessagesV2(ctx context.Context, chatID model.ChatID) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM messages_v2 WHERE chat_id = $1`, chatID.Hex()).Scan(&count)
	return count, err
}

func (r *SQLMessageRepository) scanMessageV2(row interface{ Scan(...interface{}) error }) (*model.MessageDocumentV2, error) {
	var doc model.MessageDocumentV2
	var msgIDStr, chatIDStr, contentJSON string
	var editedInt, seenInt int

	err := row.Scan(&msgIDStr, &chatIDStr, &doc.SenderID, &doc.SenderName, &doc.SenderLastName, &doc.SenderUsername,
		&doc.CreatedAt, &editedInt, &seenInt, &doc.Type, &contentJSON)
	if err != nil {
		return nil, err
	}

	doc.ID, _ = primitive.ObjectIDFromHex(msgIDStr)
	doc.ChatID, _ = primitive.ObjectIDFromHex(chatIDStr)
	doc.Edited = editedInt != 0
	doc.Seen = seenInt != 0
	json.Unmarshal([]byte(contentJSON), &doc.Content)

	return &doc, nil
}


