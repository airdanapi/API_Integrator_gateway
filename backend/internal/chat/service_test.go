package chat

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/repository"
)

// ─── stubs ────────────────────────────────────────────────────────────────────

type stubMessageStore struct {
	messages        []model.ChatMessage
	insertID        int64
	unread          int64
	unreadByConv    int64
	latest          model.ChatMessage
	latestErr       error
	insertErr       error
	listErr         error
	markReadErr     error
	countErr        error
	countByConvErr  error
	latestByConvErr error
	seenMarkConv    string
	seenMarkUser    string
}

func (s *stubMessageStore) Insert(_ context.Context, msg model.ChatMessage) (int64, error) {
	if s.insertErr != nil {
		return 0, s.insertErr
	}
	s.messages = append(s.messages, msg)
	return s.insertID, nil
}

func (s *stubMessageStore) ListByConversation(_ context.Context, _ string, limit, offset int) ([]model.ChatMessage, error) {
	if s.listErr != nil {
		return nil, s.listErr
	}
	return s.messages, nil
}

func (s *stubMessageStore) LatestByConversation(_ context.Context, _ string) (model.ChatMessage, error) {
	if s.latestByConvErr != nil {
		return model.ChatMessage{}, s.latestByConvErr
	}
	if s.latestErr != nil {
		return model.ChatMessage{}, s.latestErr
	}
	return s.latest, nil
}

func (s *stubMessageStore) MarkAsRead(_ context.Context, conversationID, toUser string) error {
	s.seenMarkConv = conversationID
	s.seenMarkUser = toUser
	return s.markReadErr
}

func (s *stubMessageStore) CountUnread(_ context.Context, _ string) (int64, error) {
	if s.countErr != nil {
		return 0, s.countErr
	}
	return s.unread, nil
}

func (s *stubMessageStore) CountUnreadByConversation(_ context.Context, _, _ string) (int64, error) {
	if s.countByConvErr != nil {
		return 0, s.countByConvErr
	}
	return s.unreadByConv, nil
}

type stubUserDirectory struct {
	users       []model.User
	admin       model.User
	adminErr    error
	findErr     error
	listErr     error
}

func (s *stubUserDirectory) FindByUsernameAndApp(_ context.Context, username, appName string) (model.User, error) {
	if s.findErr != nil {
		return model.User{}, s.findErr
	}
	for _, u := range s.users {
		if u.Username == username && u.AppName == appName {
			return u, nil
		}
	}
	return model.User{}, repository.ErrUserNotFound
}

func (s *stubUserDirectory) FindFirstByRole(_ context.Context, role model.Role) (model.User, error) {
	if s.adminErr != nil {
		return model.User{}, s.adminErr
	}
	if s.admin.Role == role {
		return s.admin, nil
	}
	return model.User{}, repository.ErrUserNotFound
}

func (s *stubUserDirectory) ListByRole(_ context.Context, role model.Role) ([]model.User, error) {
	if s.listErr != nil {
		return nil, s.listErr
	}
	var result []model.User
	for _, u := range s.users {
		if u.Role == role {
			result = append(result, u)
		}
	}
	return result, nil
}

// ─── tests ────────────────────────────────────────────────────────────────────

func fixedTime() time.Time {
	return time.Date(2026, 6, 21, 10, 0, 0, 0, time.UTC)
}

func adminClaims() auth.Claims {
	return auth.Claims{Username: "admin", Role: model.RoleAdminGateway, AppName: "API Gateway"}
}

func appUserClaims() auth.Claims {
	return auth.Claims{Username: "marketplace", Role: model.RoleAppUser, AppName: "Marketplace"}
}

func monitoringClaims() auth.Claims {
	return auth.Claims{Username: "insight", Role: model.RoleMonitoringUser, AppName: "UMKM Insight"}
}

func TestListConversations_AdminSeesAllAppUsers(t *testing.T) {
	store := &stubMessageStore{unread: 2, unreadByConv: 1}
	users := &stubUserDirectory{
		users: []model.User{
			{ID: 2, Username: "marketplace", Role: model.RoleAppUser, AppName: "Marketplace"},
			{ID: 3, Username: "pos", Role: model.RoleAppUser, AppName: "POS"},
		},
	}
	svc := New(store, users, fixedTime)

	result, err := svc.ListConversations(context.Background(), adminClaims())
	if err != nil {
		t.Fatalf("ListConversations() error: %v", err)
	}
	if len(result.Conversations) != 2 {
		t.Fatalf("conversations count = %d, want 2", len(result.Conversations))
	}
	if result.TotalUnread != 2 {
		t.Errorf("TotalUnread = %d, want 2", result.TotalUnread)
	}
	if result.Conversations[0].TargetAppName != "Marketplace" {
		t.Errorf("first conversation target = %s, want Marketplace", result.Conversations[0].TargetAppName)
	}
}

func TestListConversations_AppUserSeesAdmin(t *testing.T) {
	store := &stubMessageStore{unread: 1, unreadByConv: 1}
	users := &stubUserDirectory{
		admin: model.User{ID: 1, Username: "admin", Role: model.RoleAdminGateway, AppName: "API Gateway"},
	}
	svc := New(store, users, fixedTime)

	result, err := svc.ListConversations(context.Background(), appUserClaims())
	if err != nil {
		t.Fatalf("ListConversations() error: %v", err)
	}
	if len(result.Conversations) != 1 {
		t.Fatalf("conversations count = %d, want 1", len(result.Conversations))
	}
	if result.Conversations[0].TargetUsername != "admin" {
		t.Errorf("target = %s, want admin", result.Conversations[0].TargetUsername)
	}
}

func TestListConversations_ForbiddenForMonitoring(t *testing.T) {
	store := &stubMessageStore{}
	users := &stubUserDirectory{}
	svc := New(store, users, fixedTime)

	_, err := svc.ListConversations(context.Background(), monitoringClaims())
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("ListConversations() error = %v, want ErrForbidden", err)
	}
}

func TestSendMessage_AdminToUser(t *testing.T) {
	store := &stubMessageStore{insertID: 7, unread: 0}
	users := &stubUserDirectory{
		users: []model.User{
			{ID: 2, Username: "marketplace", Role: model.RoleAppUser, AppName: "Marketplace"},
		},
	}
	svc := New(store, users, fixedTime)

	result, err := svc.SendMessage(context.Background(), adminClaims(), SendMessageRequest{
		ToUsername: "marketplace",
		ToAppName:  "Marketplace",
		Message:    "Halo Marketplace",
	})
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}
	if result.Message.ID != 7 {
		t.Errorf("message ID = %d, want 7", result.Message.ID)
	}
	if result.Message.FromUser != "admin" {
		t.Errorf("from_user = %s, want admin", result.Message.FromUser)
	}
	if result.Message.ToUser != "marketplace" {
		t.Errorf("to_user = %s, want marketplace", result.Message.ToUser)
	}
}

func TestSendMessage_UserToAdmin(t *testing.T) {
	store := &stubMessageStore{insertID: 8, unread: 1}
	users := &stubUserDirectory{
		admin: model.User{ID: 1, Username: "admin", Role: model.RoleAdminGateway, AppName: "API Gateway"},
	}
	svc := New(store, users, fixedTime)

	result, err := svc.SendMessage(context.Background(), appUserClaims(), SendMessageRequest{
		Message: "Butuh bantuan",
	})
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}
	if result.Message.ToUser != "admin" {
		t.Errorf("to_user = %s, want admin", result.Message.ToUser)
	}
	if result.Message.FromUser != "marketplace" {
		t.Errorf("from_user = %s, want marketplace", result.Message.FromUser)
	}
}

func TestSendMessage_EmptyMessage(t *testing.T) {
	store := &stubMessageStore{}
	users := &stubUserDirectory{}
	svc := New(store, users, fixedTime)

	_, err := svc.SendMessage(context.Background(), adminClaims(), SendMessageRequest{
		ToUsername: "marketplace",
		ToAppName:  "Marketplace",
		Message:    "",
	})
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("SendMessage() error = %v, want ErrInvalidRequest", err)
	}
}

func TestSendMessage_ForbiddenRole(t *testing.T) {
	store := &stubMessageStore{}
	users := &stubUserDirectory{}
	svc := New(store, users, fixedTime)

	_, err := svc.SendMessage(context.Background(), monitoringClaims(), SendMessageRequest{
		Message: "Halo",
	})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("SendMessage() error = %v, want ErrForbidden", err)
	}
}

func TestSendMessage_TargetNotFound(t *testing.T) {
	store := &stubMessageStore{}
	users := &stubUserDirectory{
		findErr: repository.ErrUserNotFound,
	}
	svc := New(store, users, fixedTime)

	_, err := svc.SendMessage(context.Background(), adminClaims(), SendMessageRequest{
		ToUsername: "nobody",
		ToAppName:  "NoApp",
		Message:    "Halo",
	})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("SendMessage() error = %v, want ErrNotFound", err)
	}
}

func TestHistory_ValidConversation(t *testing.T) {
	ts := fixedTime()
	store := &stubMessageStore{
		messages: []model.ChatMessage{
			{ID: 1, ConversationID: "admin__marketplace__Marketplace", FromUser: "admin", ToUser: "marketplace", Message: "Halo", Timestamp: ts},
		},
		unread: 0,
	}
	users := &stubUserDirectory{
		users: []model.User{
			{ID: 2, Username: "marketplace", Role: model.RoleAppUser, AppName: "Marketplace"},
		},
	}
	svc := New(store, users, fixedTime)

	result, err := svc.History(context.Background(), adminClaims(), "admin__marketplace__Marketplace", 1, 20)
	if err != nil {
		t.Fatalf("History() error: %v", err)
	}
	if len(result.Messages) != 1 {
		t.Fatalf("messages count = %d, want 1", len(result.Messages))
	}
	if result.Page != 1 || result.Limit != 20 {
		t.Errorf("pagination = %d/%d, want 1/20", result.Page, result.Limit)
	}
}

func TestHistory_EmptyConversationID(t *testing.T) {
	store := &stubMessageStore{}
	users := &stubUserDirectory{}
	svc := New(store, users, fixedTime)

	_, err := svc.History(context.Background(), adminClaims(), "", 1, 20)
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("History() error = %v, want ErrInvalidRequest", err)
	}
}

func TestHistory_NotVisibleConversation(t *testing.T) {
	store := &stubMessageStore{}
	users := &stubUserDirectory{
		users: []model.User{},
	}
	svc := New(store, users, fixedTime)

	_, err := svc.History(context.Background(), adminClaims(), "admin__unknown__Unknown", 1, 20)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("History() error = %v, want ErrNotFound", err)
	}
}

func TestHistory_ForbiddenRole(t *testing.T) {
	store := &stubMessageStore{}
	users := &stubUserDirectory{}
	svc := New(store, users, fixedTime)

	_, err := svc.History(context.Background(), monitoringClaims(), "any__id", 1, 20)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("History() error = %v, want ErrForbidden", err)
	}
}

func TestMarkRead_Success(t *testing.T) {
	store := &stubMessageStore{unread: 0}
	users := &stubUserDirectory{
		users: []model.User{
			{ID: 2, Username: "marketplace", Role: model.RoleAppUser, AppName: "Marketplace"},
		},
	}
	svc := New(store, users, fixedTime)

	result, err := svc.MarkRead(context.Background(), adminClaims(), MarkReadRequest{
		ConversationID: "admin__marketplace__Marketplace",
	})
	if err != nil {
		t.Fatalf("MarkRead() error: %v", err)
	}
	if result.TotalUnread != 0 {
		t.Errorf("TotalUnread = %d, want 0", result.TotalUnread)
	}
	if store.seenMarkConv != "admin__marketplace__Marketplace" {
		t.Errorf("seen conversation = %s", store.seenMarkConv)
	}
	if store.seenMarkUser != "admin" {
		t.Errorf("seen user = %s", store.seenMarkUser)
	}
}

func TestMarkRead_EmptyConversationID(t *testing.T) {
	store := &stubMessageStore{}
	users := &stubUserDirectory{}
	svc := New(store, users, fixedTime)

	_, err := svc.MarkRead(context.Background(), adminClaims(), MarkReadRequest{ConversationID: ""})
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("MarkRead() error = %v, want ErrInvalidRequest", err)
	}
}

func TestMarkRead_NotFound(t *testing.T) {
	store := &stubMessageStore{}
	users := &stubUserDirectory{users: []model.User{}}
	svc := New(store, users, fixedTime)

	_, err := svc.MarkRead(context.Background(), adminClaims(), MarkReadRequest{
		ConversationID: "admin__unknown__Unknown",
	})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("MarkRead() error = %v, want ErrNotFound", err)
	}
}

func TestMarkRead_ForbiddenRole(t *testing.T) {
	store := &stubMessageStore{}
	users := &stubUserDirectory{}
	svc := New(store, users, fixedTime)

	_, err := svc.MarkRead(context.Background(), monitoringClaims(), MarkReadRequest{
		ConversationID: "any__id",
	})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("MarkRead() error = %v, want ErrForbidden", err)
	}
}

func TestService_NilDependencies(t *testing.T) {
	svc := New(nil, nil, nil)

	_, err := svc.ListConversations(context.Background(), adminClaims())
	if err == nil {
		t.Fatal("expected error for nil dependencies")
	}
}

func TestNormalizePagination(t *testing.T) {
	page, limit := normalizePagination(0, 0)
	if page != DefaultPage || limit != DefaultLimit {
		t.Errorf("got %d/%d, want %d/%d", page, limit, DefaultPage, DefaultLimit)
	}

	page, limit = normalizePagination(5, 100)
	if page != 5 || limit != DefaultLimit {
		t.Errorf("got %d/%d, want 5/%d", page, limit, DefaultLimit)
	}
}
