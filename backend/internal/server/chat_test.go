package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/config"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/chat"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/server"
)

type stubChatSvc struct {
	conversations chat.ConversationListResult
	history       chat.HistoryResult
	sendResult    chat.SendMessageResult
	readResult    chat.MarkReadResult
	err           error
	seenClaims    auth.Claims
	seenHistoryID string
	seenPage      int
	seenLimit     int
	seenSend      chat.SendMessageRequest
	seenRead      chat.MarkReadRequest
}

func (s *stubChatSvc) ListConversations(_ context.Context, claims auth.Claims) (chat.ConversationListResult, error) {
	s.seenClaims = claims
	return s.conversations, s.err
}

func (s *stubChatSvc) History(_ context.Context, claims auth.Claims, conversationID string, page, limit int) (chat.HistoryResult, error) {
	s.seenClaims = claims
	s.seenHistoryID = conversationID
	s.seenPage = page
	s.seenLimit = limit
	return s.history, s.err
}

func (s *stubChatSvc) SendMessage(_ context.Context, claims auth.Claims, req chat.SendMessageRequest) (chat.SendMessageResult, error) {
	s.seenClaims = claims
	s.seenSend = req
	return s.sendResult, s.err
}

func (s *stubChatSvc) MarkRead(_ context.Context, claims auth.Claims, req chat.MarkReadRequest) (chat.MarkReadResult, error) {
	s.seenClaims = claims
	s.seenRead = req
	return s.readResult, s.err
}

func TestChatRequiresAuthentication(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{ChatService: &stubChatSvc{}})
	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/chat/conversations", nil))
	if err != nil {
		t.Fatalf("GET /chat/conversations error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", resp.StatusCode)
	}
}

func TestChatConversationsSuccess(t *testing.T) {
	svc := &stubChatSvc{conversations: chat.ConversationListResult{
		Conversations: []chat.Conversation{{ConversationID: "admin__marketplace__Marketplace", TargetUsername: "marketplace", TargetAppName: "Marketplace", UnreadCount: 1}},
		TotalUnread:   1,
	}}
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier: makeChatVerifier(model.RoleAdminGateway, "admin", "API Gateway"),
		ChatService:   svc,
	})
	req := httptest.NewRequest(http.MethodGet, "/chat/conversations", nil)
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("GET /chat/conversations error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if svc.seenClaims.Username != "admin" {
		t.Fatalf("seen claims = %#v", svc.seenClaims)
	}
	var body struct {
		Status string                      `json:"status"`
		Data   chat.ConversationListResult `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Status != "success" || body.Data.TotalUnread != 1 || len(body.Data.Conversations) != 1 {
		t.Fatalf("body = %#v", body)
	}
}

func TestChatHistoryValidatesConversationAndMapsNotFound(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier: makeChatVerifier(model.RoleAppUser, "marketplace", "Marketplace"),
		ChatService:   &stubChatSvc{},
	})
	req := httptest.NewRequest(http.MethodGet, "/chat/history", nil)
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("GET /chat/history error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("missing conversation status = %d, want 400", resp.StatusCode)
	}

	app = server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier: makeChatVerifier(model.RoleAppUser, "marketplace", "Marketplace"),
		ChatService:   &stubChatSvc{err: chat.ErrNotFound},
	})
	req = httptest.NewRequest(http.MethodGet, "/chat/history?conversation_id=admin__pos__POS", nil)
	req.Header.Set("Authorization", "Bearer token")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("GET /chat/history not found error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("not found status = %d, want 404", resp.StatusCode)
	}
}

func TestChatHistorySuccess(t *testing.T) {
	svc := &stubChatSvc{history: chat.HistoryResult{
		Messages: []chat.Message{{ID: 1, ConversationID: "admin__marketplace__Marketplace", FromUser: "admin", ToUser: "marketplace", Message: "Halo", Timestamp: time.Date(2026, 6, 21, 10, 0, 0, 0, time.UTC)}},
		Page:     1,
		Limit:    20,
	}}
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier: makeChatVerifier(model.RoleAppUser, "marketplace", "Marketplace"),
		ChatService:   svc,
	})
	req := httptest.NewRequest(http.MethodGet, "/chat/history?conversation_id=admin__marketplace__Marketplace&page=1&limit=20", nil)
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("GET /chat/history error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if svc.seenHistoryID != "admin__marketplace__Marketplace" || svc.seenPage != 1 || svc.seenLimit != 20 {
		t.Fatalf("seen args id=%q page=%d limit=%d", svc.seenHistoryID, svc.seenPage, svc.seenLimit)
	}
}

func TestChatSendValidatesBodyAndMapsErrors(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier: makeChatVerifier(model.RoleAdminGateway, "admin", "API Gateway"),
		ChatService:   &stubChatSvc{},
	})
	req := httptest.NewRequest(http.MethodPost, "/chat/message", bytes.NewBufferString(`{}`))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /chat/message error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid body status = %d, want 400", resp.StatusCode)
	}

	app = server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier: makeChatVerifier(model.RoleMonitoringUser, "insight", "UMKM Insight"),
		ChatService:   &stubChatSvc{err: chat.ErrForbidden},
	})
	req = httptest.NewRequest(http.MethodPost, "/chat/message", bytes.NewBufferString(`{"message":"Halo"}`))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("POST /chat/message forbidden error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("forbidden status = %d, want 403", resp.StatusCode)
	}
}

func TestChatSendAndReadSuccess(t *testing.T) {
	svc := &stubChatSvc{
		sendResult: chat.SendMessageResult{Message: chat.Message{ID: 7, ConversationID: "admin__marketplace__Marketplace", Message: "Halo"}},
		readResult: chat.MarkReadResult{TotalUnread: 0},
	}
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier: makeChatVerifier(model.RoleAdminGateway, "admin", "API Gateway"),
		ChatService:   svc,
	})
	req := httptest.NewRequest(http.MethodPost, "/chat/message", bytes.NewBufferString(`{"to_username":"marketplace","to_app_name":"Marketplace","message":"Halo"}`))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /chat/message error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("send status = %d, want 200", resp.StatusCode)
	}
	if svc.seenSend.ToUsername != "marketplace" || svc.seenSend.Message != "Halo" {
		t.Fatalf("seenSend = %#v", svc.seenSend)
	}

	req = httptest.NewRequest(http.MethodPost, "/chat/read", bytes.NewBufferString(`{"conversation_id":"admin__marketplace__Marketplace"}`))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("POST /chat/read error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || svc.seenRead.ConversationID != "admin__marketplace__Marketplace" {
		t.Fatalf("read status = %d seen=%#v", resp.StatusCode, svc.seenRead)
	}
}

func makeChatVerifier(role model.Role, username, appName string) *stubVerifier {
	return &stubVerifier{claims: auth.Claims{Username: username, Role: role, AppName: appName}}
}
