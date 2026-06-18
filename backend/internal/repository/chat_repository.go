package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

// ChatRepository mendefinisikan kontrak akses data chat_messages.
type ChatRepository interface {
	Insert(ctx context.Context, msg model.ChatMessage) (int64, error)
	ListByConversation(ctx context.Context, conversationID string, limit, offset int) ([]model.ChatMessage, error)
	ListConversations(ctx context.Context, username string) ([]string, error)
	MarkAsRead(ctx context.Context, conversationID, toUser string) error
	CountUnread(ctx context.Context, toUser string) (int64, error)
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
		 ORDER BY timestamp ASC LIMIT ? OFFSET ?`,
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
		`SELECT DISTINCT conversation_id FROM chat_messages
		 WHERE from_user = ? OR to_user = ?
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

// scanChatMessages mem-parse baris SQL menjadi slice ChatMessage.
func scanChatMessages(rows *sql.Rows) ([]model.ChatMessage, error) {
	var msgs []model.ChatMessage
	for rows.Next() {
		var m model.ChatMessage
		var isRead int
		if err := rows.Scan(
			&m.ID, &m.ConversationID, &m.FromUser, &m.ToUser,
			&m.Message, &m.Timestamp, &isRead,
		); err != nil {
			return nil, fmt.Errorf("scan chat message row: %w", err)
		}
		m.IsRead = isRead == 1
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}
