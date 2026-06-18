package model

import "time"

// NotificationType adalah jenis-jenis notifikasi yang didukung sistem.
type NotificationType string

const (
	NotificationTypeAPIInactive   NotificationType = "api_inactive"
	NotificationTypeErrorRate     NotificationType = "error_rate"
	NotificationTypeResponseTime  NotificationType = "response_time"
	NotificationTypeSystem        NotificationType = "system"
)

// RequestLog merepresentasikan satu entri audit log request gateway.
type RequestLog struct {
	ID         int64
	Timestamp  time.Time
	SourceApp  string
	Endpoint   string
	Method     string
	Payload    []byte // JSON raw
	Status     int
	Response   []byte // JSON raw
	DurationMS *int
}

// Notification merepresentasikan satu notifikasi sistem.
type Notification struct {
	ID        int64
	CreatedAt time.Time
	AppName   string
	Type      NotificationType
	Message   string
	IsRead    bool
}

// ChatMessage merepresentasikan satu pesan dalam percakapan chat.
type ChatMessage struct {
	ID             int64
	ConversationID string
	FromUser       string
	ToUser         string
	Message        string
	Timestamp      time.Time
	IsRead         bool
}

// DashboardData merepresentasikan satu entri cache analytics dashboard.
type DashboardData struct {
	ID         int64
	CacheKey   string
	AppName    string
	Data       []byte // JSON raw
	ComputedAt time.Time
	ExpiresAt  time.Time
}
