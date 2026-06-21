package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

var ErrChatMessageNotFound = errors.New("chat message not found")

// ChatRepository mendefinisikan kontrak akses data chat_messages.
type ChatRepository interface {
	Insert(ctx context.Context, msg model.ChatMessage) (int64, error)
	ListByConversation(ctx context.Context, conversationID string, limit, offset int) ([]model.ChatMessage, error)
	ListConversations(ctx context.Context, username string) ([]string, error)
	LatestByConversation(ctx context.Context, conversationID string) (model.ChatMessage, error)
	MarkAsRead(ctx context.Context, conversationID, toUser string) error
	CountUnread(ctx context.Context, toUser string) (int64, error)
	CountUnreadByConversation(ctx context.Context, conversationID, toUser string) (int64, error)
}

// MySQLChatRepository mengimplementasikan ChatRepository dengan MySQL.
type MySQLChatRepository struct {
	db *sql.DB
}

// NewMySQLChatRepository membuat instance baru.
func NewMySQLChatRepository(db *sql.DB) *MySQLChatRepository {
	return &MySQLChatRepository{db: db}
}

// Insert menyimpan satu pesan chat dan mengembalikan ID-nya.
func (r *MySQLChatRepository) Insert(ctx context.Context, msg model.ChatMessage) (int64, error) {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO chat_messages (conversation_id, from_user, to_user, message, timestamp, is_read)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		msg.ConversationID, msg.FromUser, msg.ToUser, msg.Message, msg.Timestamp, boolToInt(msg.IsRead),
	)
	if err != nil {
		return 0, fmt.Errorf("insert chat message: %w", err)
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// ListByConversation mengambil pesan dalam satu conversation, terlama dulu (natural chat order).
func (r *MySQLChatRepository) ListByConversation(ctx context.Context, conversationID string, limit, offset int) ([]model.ChatMessage, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, conversation_id, from_user, to_user, message, timestamp, is_read
		 FROM chat_messages WHERE conversation_id = ?
		 ORDER BY timestamp ASC, id ASC LIMIT ? OFFSET ?`,
		conversationID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list chat messages by conversation: %w", err)
	}
	defer rows.Close()
	return scanChatMessages(rows)
}

// ListConversations mengambil daftar conversation_id unik yang melibatkan user tertentu.
func (r *MySQLChatRepository) ListConversations(ctx context.Context, username string) ([]string, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT conversation_id FROM chat_messages
		 WHERE from_user = ? OR to_user = ?
		 GROUP BY conversation_id
		 ORDER BY MAX(timestamp) DESC`,
		username, username,
	)
	if err != nil {
		return nil, fmt.Errorf("list conversations for %s: %w", username, err)
	}
	defer rows.Close()
	var convIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan conversation id: %w", err)
		}
		convIDs = append(convIDs, id)
	}
	return convIDs, rows.Err()
}

// LatestByConversation mengambil pesan terbaru dalam satu conversation.
func (r *MySQLChatRepository) LatestByConversation(ctx context.Context, conversationID string) (model.ChatMessage, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, conversation_id, from_user, to_user, message, timestamp, is_read
		 FROM chat_messages WHERE conversation_id = ?
		 ORDER BY timestamp DESC, id DESC LIMIT 1`,
		conversationID,
	)
	msg, err := scanChatMessage(row)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ChatMessage{}, ErrChatMessageNotFound
	}
	if err != nil {
		return model.ChatMessage{}, fmt.Errorf("latest chat message for %s: %w", conversationID, err)
	}
	return msg, nil
}

// MarkAsRead menandai semua pesan dalam conversation yang ditujukan ke toUser sebagai dibaca.
func (r *MySQLChatRepository) MarkAsRead(ctx context.Context, conversationID, toUser string) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE chat_messages SET is_read = 1
		 WHERE conversation_id = ? AND to_user = ? AND is_read = 0`,
		conversationID, toUser,
	)
	if err != nil {
		return fmt.Errorf("mark chat messages as read in conversation %s for %s: %w", conversationID, toUser, err)
	}
	return nil
}

// CountUnread menghitung total pesan yang belum dibaca untuk penerima tertentu.
func (r *MySQLChatRepository) CountUnread(ctx context.Context, toUser string) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM chat_messages WHERE to_user = ? AND is_read = 0`,
		toUser,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count unread chat messages for %s: %w", toUser, err)
	}
	return count, nil
}

// CountUnreadByConversation menghitung pesan belum dibaca dalam satu conversation untuk penerima tertentu.
func (r *MySQLChatRepository) CountUnreadByConversation(ctx context.Context, conversationID, toUser string) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM chat_messages
		 WHERE conversation_id = ? AND to_user = ? AND is_read = 0`,
		conversationID, toUser,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count unread chat messages for %s/%s: %w", conversationID, toUser, err)
	}
	return count, nil
}

// scanChatMessages mem-parse baris SQL menjadi slice ChatMessage.
func scanChatMessages(rows *sql.Rows) ([]model.ChatMessage, error) {
	var msgs []model.ChatMessage
	for rows.Next() {
		msg, err := scanChatMessage(rows)
		if err != nil {
			return nil, fmt.Errorf("scan chat message row: %w", err)
		}
		msgs = append(msgs, msg)
	}
	return msgs, rows.Err()
}

type chatMessageScanner interface {
	Scan(dest ...any) error
}

func scanChatMessage(row chatMessageScanner) (model.ChatMessage, error) {
	var m model.ChatMessage
	var isRead int
	if err := row.Scan(
		&m.ID, &m.ConversationID, &m.FromUser, &m.ToUser,
		&m.Message, &m.Timestamp, &isRead,
	); err != nil {
		return model.ChatMessage{}, err
	}
	m.IsRead = isRead == 1
	return m, nil
}
