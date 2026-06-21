package chat

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/repository"
)

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 50

	maxMessageLength = 1000
)

var (
	ErrInvalidRequest = errors.New("invalid chat request")
	ErrForbidden      = errors.New("chat access forbidden")
	ErrNotFound       = errors.New("chat conversation not found")
)

type MessageStore interface {
	Insert(ctx context.Context, msg model.ChatMessage) (int64, error)
	ListByConversation(ctx context.Context, conversationID string, limit, offset int) ([]model.ChatMessage, error)
	LatestByConversation(ctx context.Context, conversationID string) (model.ChatMessage, error)
	MarkAsRead(ctx context.Context, conversationID, toUser string) error
	CountUnread(ctx context.Context, toUser string) (int64, error)
	CountUnreadByConversation(ctx context.Context, conversationID, toUser string) (int64, error)
}

type UserDirectory interface {
	FindByUsernameAndApp(ctx context.Context, username string, appName string) (model.User, error)
	FindFirstByRole(ctx context.Context, role model.Role) (model.User, error)
	ListByRole(ctx context.Context, role model.Role) ([]model.User, error)
}

type Conversation struct {
	ConversationID string   `json:"conversation_id"`
	TargetUsername string   `json:"target_username"`
	TargetAppName  string   `json:"target_app_name"`
	UnreadCount    int64    `json:"unread_count"`
	LatestMessage  *Message `json:"latest_message,omitempty"`
}

type ConversationListResult struct {
	Conversations []Conversation `json:"conversations"`
	TotalUnread   int64          `json:"total_unread"`
}

type Message struct {
	ID             int64     `json:"id"`
	ConversationID string    `json:"conversation_id"`
	FromUser       string    `json:"from_user"`
	ToUser         string    `json:"to_user"`
	Message        string    `json:"message"`
	Timestamp      time.Time `json:"timestamp"`
	IsRead         bool      `json:"is_read"`
}

type HistoryResult struct {
	Messages    []Message `json:"messages"`
	TotalUnread int64     `json:"total_unread"`
	Page        int       `json:"page"`
	Limit       int       `json:"limit"`
}

type SendMessageRequest struct {
	ToUsername string `json:"to_username"`
	ToAppName  string `json:"to_app_name"`
	Message    string `json:"message"`
}

type SendMessageResult struct {
	Message     Message `json:"message"`
	TotalUnread int64   `json:"total_unread"`
}

type MarkReadRequest struct {
	ConversationID string `json:"conversation_id"`
}

type MarkReadResult struct {
	TotalUnread int64 `json:"total_unread"`
}

type Service struct {
	messages MessageStore
	users    UserDirectory
	now      func() time.Time
}

func New(messages MessageStore, users UserDirectory, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{messages: messages, users: users, now: now}
}

func (s *Service) ListConversations(ctx context.Context, claims auth.Claims) (ConversationListResult, error) {
	if err := s.ready(); err != nil {
		return ConversationListResult{}, err
	}
	switch claims.Role {
	case model.RoleAdminGateway:
		appUsers, err := s.users.ListByRole(ctx, model.RoleAppUser)
		if err != nil {
			return ConversationListResult{}, err
		}
		conversations := make([]Conversation, 0, len(appUsers))
		for _, appUser := range appUsers {
			conversationID := conversationID(claims.Username, appUser.Username, appUser.AppName)
			conversation, err := s.conversationForTarget(ctx, conversationID, claims.Username, appUser.Username, appUser.AppName)
			if err != nil {
				return ConversationListResult{}, err
			}
			conversations = append(conversations, conversation)
		}
		totalUnread, err := s.messages.CountUnread(ctx, claims.Username)
		if err != nil {
			return ConversationListResult{}, err
		}
		return ConversationListResult{Conversations: conversations, TotalUnread: totalUnread}, nil
	case model.RoleAppUser:
		admin, err := s.users.FindFirstByRole(ctx, model.RoleAdminGateway)
		if errors.Is(err, repository.ErrUserNotFound) {
			return ConversationListResult{}, ErrNotFound
		}
		if err != nil {
			return ConversationListResult{}, err
		}
		conversationID := conversationID(admin.Username, claims.Username, claims.AppName)
		conversation, err := s.conversationForTarget(ctx, conversationID, claims.Username, admin.Username, admin.AppName)
		if err != nil {
			return ConversationListResult{}, err
		}
		totalUnread, err := s.messages.CountUnread(ctx, claims.Username)
		if err != nil {
			return ConversationListResult{}, err
		}
		return ConversationListResult{Conversations: []Conversation{conversation}, TotalUnread: totalUnread}, nil
	default:
		return ConversationListResult{}, ErrForbidden
	}
}

func (s *Service) History(ctx context.Context, claims auth.Claims, id string, page, limit int) (HistoryResult, error) {
	if err := s.ready(); err != nil {
		return HistoryResult{}, err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return HistoryResult{}, ErrInvalidRequest
	}
	if ok, err := s.visibleConversation(ctx, claims, id); err != nil {
		return HistoryResult{}, err
	} else if !ok {
		return HistoryResult{}, ErrNotFound
	}
	page, limit = normalizePagination(page, limit)
	offset := (page - 1) * limit
	items, err := s.messages.ListByConversation(ctx, id, limit, offset)
	if err != nil {
		return HistoryResult{}, err
	}
	totalUnread, err := s.messages.CountUnread(ctx, claims.Username)
	if err != nil {
		return HistoryResult{}, err
	}
	return HistoryResult{
		Messages:    toMessages(items),
		TotalUnread: totalUnread,
		Page:        page,
		Limit:       limit,
	}, nil
}

func (s *Service) SendMessage(ctx context.Context, claims auth.Claims, req SendMessageRequest) (SendMessageResult, error) {
	if err := s.ready(); err != nil {
		return SendMessageResult{}, err
	}
	messageText := strings.TrimSpace(req.Message)
	if messageText == "" || len(messageText) > maxMessageLength {
		return SendMessageResult{}, ErrInvalidRequest
	}

	var (
		id     string
		toUser string
	)
	switch claims.Role {
	case model.RoleAdminGateway:
		targetUsername := strings.TrimSpace(req.ToUsername)
		targetAppName := strings.TrimSpace(req.ToAppName)
		if targetUsername == "" || targetAppName == "" {
			return SendMessageResult{}, ErrInvalidRequest
		}
		target, err := s.users.FindByUsernameAndApp(ctx, targetUsername, targetAppName)
		if errors.Is(err, repository.ErrUserNotFound) {
			return SendMessageResult{}, ErrNotFound
		}
		if err != nil {
			return SendMessageResult{}, err
		}
		if target.Role != model.RoleAppUser {
			return SendMessageResult{}, ErrNotFound
		}
		id = conversationID(claims.Username, target.Username, target.AppName)
		toUser = target.Username
	case model.RoleAppUser:
		admin, err := s.users.FindFirstByRole(ctx, model.RoleAdminGateway)
		if errors.Is(err, repository.ErrUserNotFound) {
			return SendMessageResult{}, ErrNotFound
		}
		if err != nil {
			return SendMessageResult{}, err
		}
		id = conversationID(admin.Username, claims.Username, claims.AppName)
		toUser = admin.Username
	default:
		return SendMessageResult{}, ErrForbidden
	}

	msg := model.ChatMessage{
		ConversationID: id,
		FromUser:       claims.Username,
		ToUser:         toUser,
		Message:        messageText,
		Timestamp:      s.now().UTC(),
		IsRead:         false,
	}
	insertedID, err := s.messages.Insert(ctx, msg)
	if err != nil {
		return SendMessageResult{}, err
	}
	msg.ID = insertedID
	totalUnread, err := s.messages.CountUnread(ctx, claims.Username)
	if err != nil {
		return SendMessageResult{}, err
	}
	return SendMessageResult{Message: toMessage(msg), TotalUnread: totalUnread}, nil
}

func (s *Service) MarkRead(ctx context.Context, claims auth.Claims, req MarkReadRequest) (MarkReadResult, error) {
	if err := s.ready(); err != nil {
		return MarkReadResult{}, err
	}
	id := strings.TrimSpace(req.ConversationID)
	if id == "" {
		return MarkReadResult{}, ErrInvalidRequest
	}
	if ok, err := s.visibleConversation(ctx, claims, id); err != nil {
		return MarkReadResult{}, err
	} else if !ok {
		return MarkReadResult{}, ErrNotFound
	}
	if err := s.messages.MarkAsRead(ctx, id, claims.Username); err != nil {
		return MarkReadResult{}, err
	}
	totalUnread, err := s.messages.CountUnread(ctx, claims.Username)
	if err != nil {
		return MarkReadResult{}, err
	}
	return MarkReadResult{TotalUnread: totalUnread}, nil
}

func (s *Service) conversationForTarget(ctx context.Context, id, viewerUsername, targetUsername, targetAppName string) (Conversation, error) {
	unread, err := s.messages.CountUnreadByConversation(ctx, id, viewerUsername)
	if err != nil {
		return Conversation{}, err
	}
	var latest *Message
	latestModel, err := s.messages.LatestByConversation(ctx, id)
	if err == nil {
		msg := toMessage(latestModel)
		latest = &msg
	} else if !errors.Is(err, repository.ErrChatMessageNotFound) {
		return Conversation{}, err
	}
	return Conversation{
		ConversationID: id,
		TargetUsername: targetUsername,
		TargetAppName:  targetAppName,
		UnreadCount:    unread,
		LatestMessage:  latest,
	}, nil
}

func (s *Service) visibleConversation(ctx context.Context, claims auth.Claims, id string) (bool, error) {
	switch claims.Role {
	case model.RoleAdminGateway:
		appUsers, err := s.users.ListByRole(ctx, model.RoleAppUser)
		if err != nil {
			return false, err
		}
		for _, appUser := range appUsers {
			if conversationID(claims.Username, appUser.Username, appUser.AppName) == id {
				return true, nil
			}
		}
		return false, nil
	case model.RoleAppUser:
		admin, err := s.users.FindFirstByRole(ctx, model.RoleAdminGateway)
		if errors.Is(err, repository.ErrUserNotFound) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		return conversationID(admin.Username, claims.Username, claims.AppName) == id, nil
	default:
		return false, ErrForbidden
	}
}

func (s *Service) ready() error {
	if s.messages == nil || s.users == nil {
		return errors.New("chat dependencies are nil")
	}
	return nil
}

func conversationID(adminUsername, appUsername, appName string) string {
	return strings.Join([]string{adminUsername, appUsername, appName}, "__")
}

func normalizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = DefaultPage
	}
	if limit < 1 || limit > MaxLimit {
		limit = DefaultLimit
	}
	return page, limit
}

func toMessages(items []model.ChatMessage) []Message {
	result := make([]Message, 0, len(items))
	for _, item := range items {
		result = append(result, toMessage(item))
	}
	return result
}

func toMessage(item model.ChatMessage) Message {
	return Message{
		ID:             item.ID,
		ConversationID: item.ConversationID,
		FromUser:       item.FromUser,
		ToUser:         item.ToUser,
		Message:        item.Message,
		Timestamp:      item.Timestamp,
		IsRead:         item.IsRead,
	}
}
