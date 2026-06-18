package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

func TestChatRepository_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	msg := model.ChatMessage{
		ConversationID: "conv-admin-marketplace",
		FromUser:       "admin",
		ToUser:         "marketplace",
		Message:        "Halo, ada masalah?",
		Timestamp:      time.Now().UTC(),
		IsRead:         false,
	}

	mock.ExpectExec(`INSERT INTO chat_messages`).
		WithArgs(
			msg.ConversationID,
			msg.FromUser,
			msg.ToUser,
			msg.Message,
			sqlmock.AnyArg(), // timestamp
			0,                // is_read = false
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewMySQLChatRepository(db)
	id, err := repo.Insert(context.Background(), msg)
	if err != nil {
		t.Fatalf("Insert() error: %v", err)
	}
	if id != 1 {
		t.Errorf("Insert() id = %d, want 1", id)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}

func TestChatRepository_CountUnread(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM chat_messages WHERE to_user = \? AND is_read = 0`).
		WithArgs("marketplace").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	repo := NewMySQLChatRepository(db)
	count, err := repo.CountUnread(context.Background(), "marketplace")
	if err != nil {
		t.Fatalf("CountUnread() error: %v", err)
	}
	if count != 2 {
		t.Errorf("CountUnread() = %d, want 2", count)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}

func TestChatRepository_ListByConversation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	ts := time.Now().UTC()
	columns := []string{"id", "conversation_id", "from_user", "to_user", "message", "timestamp", "is_read"}
	rows := sqlmock.NewRows(columns).
		AddRow(1, "conv-admin-marketplace", "admin", "marketplace", "Halo!", ts, 1).
		AddRow(2, "conv-admin-marketplace", "marketplace", "admin", "Halo juga!", ts, 0)

	mock.ExpectQuery(`SELECT id, conversation_id, from_user, to_user, message, timestamp, is_read FROM chat_messages WHERE conversation_id = \?`).
		WithArgs("conv-admin-marketplace", 20, 0).
		WillReturnRows(rows)

	repo := NewMySQLChatRepository(db)
	msgs, err := repo.ListByConversation(context.Background(), "conv-admin-marketplace", 20, 0)
	if err != nil {
		t.Fatalf("ListByConversation() error: %v", err)
	}
	if len(msgs) != 2 {
		t.Errorf("ListByConversation() count = %d, want 2", len(msgs))
	}
	if !msgs[0].IsRead {
		t.Errorf("msgs[0].IsRead = false, want true")
	}
	if msgs[1].IsRead {
		t.Errorf("msgs[1].IsRead = true, want false")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}
