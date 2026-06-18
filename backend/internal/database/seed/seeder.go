// Package seed menyediakan fungsi untuk mengisi data awal (seeding) ke semua
// tabel Sprint 5: request_logs, notifications, chat_messages, dan dashboard_data.
//
// Seeder ini dirancang untuk idempoten: aman dijalankan berulang kali.
// Data yang sudah ada tidak akan diduplikasi (menggunakan INSERT IGNORE atau ON DUPLICATE KEY).
package seed

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/repository"
)

// Options mengonfigurasi perilaku seeder.
type Options struct {
	// Verbose mengaktifkan logging detail per record yang di-seed.
	Verbose bool
	// RequestLogCount jumlah request log yang akan di-seed (default: 50).
	RequestLogCount int
	// ChatMessageCount jumlah pesan chat yang akan di-seed (default: 20).
	ChatMessageCount int
}

// DefaultOptions mengembalikan konfigurasi seeder default.
func DefaultOptions() Options {
	return Options{
		Verbose:          false,
		RequestLogCount:  50,
		ChatMessageCount: 20,
	}
}

// Run menjalankan seeding semua tabel Sprint 5.
// Urutan: request_logs → notifications → chat_messages → dashboard_data.
func Run(ctx context.Context, db *sql.DB, opts Options) error {
	logRepo := repository.NewMySQLLogRepository(db)
	notifRepo := repository.NewMySQLNotificationRepository(db)
	chatRepo := repository.NewMySQLChatRepository(db)
	dashRepo := repository.NewMySQLDashboardRepository(db)

	if opts.RequestLogCount <= 0 {
		opts.RequestLogCount = 50
	}
	if opts.ChatMessageCount <= 0 {
		opts.ChatMessageCount = 20
	}

	steps := []struct {
		name string
		fn   func() error
	}{
		{"request_logs", func() error { return seedRequestLogs(ctx, logRepo, opts) }},
		{"notifications", func() error { return seedNotifications(ctx, notifRepo, opts) }},
		{"chat_messages", func() error { return seedChatMessages(ctx, chatRepo, opts) }},
		{"dashboard_data", func() error { return seedDashboardData(ctx, dashRepo, opts) }},
	}

	for _, step := range steps {
		if opts.Verbose {
			log.Printf("[seed] mulai seeding tabel: %s", step.name)
		}
		if err := step.fn(); err != nil {
			return fmt.Errorf("seed %s: %w", step.name, err)
		}
		if opts.Verbose {
			log.Printf("[seed] selesai seeding tabel: %s", step.name)
		}
	}
	return nil
}

// ─── request_logs ──────────────────────────────────────────────────────────────

var sourceApps = []string{
	"Marketplace", "POS", "SupplierHub", "LogistiKita", "SmartBank",
}

var endpoints = []string{
	"/gateway/payment",
	"/gateway/smartbank",
	"/gateway/marketplace",
	"/gateway/logistics",
	"/gateway/supplier",
}

var httpStatuses = []int{200, 200, 200, 200, 201, 400, 401, 500}

func seedRequestLogs(ctx context.Context, repo *repository.MySQLLogRepository, opts Options) error {
	now := time.Now().UTC()
	for i := range opts.RequestLogCount {
		app := sourceApps[i%len(sourceApps)]
		endpoint := endpoints[i%len(endpoints)]
		status := httpStatuses[i%len(httpStatuses)]
		daysAgo := i % 14 // antara 0-13 hari lalu
		ts := now.Add(-time.Duration(daysAgo) * 24 * time.Hour).
			Add(-time.Duration(i) * time.Minute)

		payload, _ := json.Marshal(map[string]any{
			"from_app":     app,
			"from_user":    fmt.Sprintf("user_%s_%d", app, i),
			"amount":       (i + 1) * 10000,
			"service_type": "payment",
		})
		response, _ := json.Marshal(map[string]any{
			"status":         statusText(status),
			"transaction_id": fmt.Sprintf("TXN-%04d", i+1),
		})
		durationMS := 50 + (i * 7 % 150) // 50–200ms simulasi

		_, err := repo.Insert(ctx, model.RequestLog{
			Timestamp:  ts,
			SourceApp:  app,
			Endpoint:   endpoint,
			Method:     "POST",
			Payload:    payload,
			Status:     status,
			Response:   response,
			DurationMS: &durationMS,
		})
		if err != nil {
			return fmt.Errorf("insert request log #%d: %w", i+1, err)
		}
		if opts.Verbose {
			log.Printf("  [request_logs] #%d: %s → %s (%d)", i+1, app, endpoint, status)
		}
	}
	return nil
}

// ─── notifications ─────────────────────────────────────────────────────────────

type notifSeed struct {
	AppName string
	Type    model.NotificationType
	Message string
	IsRead  bool
	DaysAgo int
}

var notifSeeds = []notifSeed{
	{
		AppName: "Marketplace",
		Type:    model.NotificationTypeAPIInactive,
		Message: "API Marketplace tidak aktif selama lebih dari 1 minggu.",
		IsRead:  false,
		DaysAgo: 0,
	},
	{
		AppName: "POS",
		Type:    model.NotificationTypeErrorRate,
		Message: "Error rate endpoint /gateway/payment melewati 10% dalam 1 jam terakhir.",
		IsRead:  false,
		DaysAgo: 0,
	},
	{
		AppName: "SupplierHub",
		Type:    model.NotificationTypeResponseTime,
		Message: "Rata-rata response time /gateway/supplier melebihi 300ms.",
		IsRead:  true,
		DaysAgo: 2,
	},
	{
		AppName: "LogistiKita",
		Type:    model.NotificationTypeAPIInactive,
		Message: "API LogistiKita tidak aktif selama lebih dari 1 minggu.",
		IsRead:  true,
		DaysAgo: 3,
	},
	{
		AppName: "API Gateway",
		Type:    model.NotificationTypeSystem,
		Message: "Sistem berhasil dimulai ulang setelah maintenance.",
		IsRead:  true,
		DaysAgo: 5,
	},
	{
		AppName: "SmartBank",
		Type:    model.NotificationTypeResponseTime,
		Message: "Response time SmartBank meningkat signifikan, periksa infrastruktur.",
		IsRead:  false,
		DaysAgo: 1,
	},
	{
		AppName: "UMKM Insight",
		Type:    model.NotificationTypeSystem,
		Message: "Data dashboard berhasil diperbarui.",
		IsRead:  true,
		DaysAgo: 7,
	},
}

func seedNotifications(ctx context.Context, repo *repository.MySQLNotificationRepository, opts Options) error {
	now := time.Now().UTC()
	for i, seed := range notifSeeds {
		_, err := repo.Insert(ctx, model.Notification{
			CreatedAt: now.Add(-time.Duration(seed.DaysAgo) * 24 * time.Hour),
			AppName:   seed.AppName,
			Type:      seed.Type,
			Message:   seed.Message,
			IsRead:    seed.IsRead,
		})
		if err != nil {
			return fmt.Errorf("insert notification #%d: %w", i+1, err)
		}
		if opts.Verbose {
			log.Printf("  [notifications] #%d: [%s] %s — %s", i+1, seed.Type, seed.AppName, seed.Message[:30]+"...")
		}
	}
	return nil
}

// ─── chat_messages ─────────────────────────────────────────────────────────────

type chatSeed struct {
	ConvID   string
	From     string
	To       string
	Msg      string
	MinutesAgo int
	IsRead   bool
}

func buildChatSeeds(total int) []chatSeed {
	seeds := []chatSeed{
		// Percakapan admin ↔ marketplace
		{"conv-admin-marketplace", "admin", "marketplace", "Halo, ada masalah dengan integrasi payment Marketplace?", 60, true},
		{"conv-admin-marketplace", "marketplace", "admin", "Ya, endpoint /gateway/payment kadang timeout.", 55, true},
		{"conv-admin-marketplace", "admin", "marketplace", "Kami sedang investigasi, mohon tunggu.", 50, true},
		{"conv-admin-marketplace", "marketplace", "admin", "Baik, terima kasih admin.", 45, true},
		{"conv-admin-marketplace", "admin", "marketplace", "Masalah sudah diperbaiki, silakan coba lagi.", 30, true},
		// Percakapan admin ↔ pos
		{"conv-admin-pos", "pos", "admin", "Error rate meningkat, apakah ada gangguan?", 120, true},
		{"conv-admin-pos", "admin", "pos", "Ada maintenance terjadwal kemarin, sekarang sudah normal.", 115, true},
		{"conv-admin-pos", "pos", "admin", "Terima kasih infonya!", 110, true},
		// Percakapan admin ↔ supplierhub
		{"conv-admin-supplierhub", "supplierhub", "admin", "Request supplier order sering gagal, tolong dicek.", 200, false},
		{"conv-admin-supplierhub", "admin", "supplierhub", "Sudah kami cek, ada bug di validasi payload. Akan diperbaiki.", 190, false},
		// Percakapan admin ↔ insight
		{"conv-admin-insight", "insight", "admin", "Apakah saya bisa mengakses data analytics bulan lalu?", 300, true},
		{"conv-admin-insight", "admin", "insight", "Bisa, dashboard monitoring sudah tersedia.", 295, true},
	}
	if total < len(seeds) {
		return seeds[:total]
	}
	// Tambah pesan generik jika dibutuhkan lebih banyak
	for i := len(seeds); i < total; i++ {
		pair := i % 4
		var from, to, convID string
		switch pair {
		case 0:
			from, to, convID = "admin", "marketplace", "conv-admin-marketplace"
		case 1:
			from, to, convID = "admin", "pos", "conv-admin-pos"
		case 2:
			from, to, convID = "admin", "supplierhub", "conv-admin-supplierhub"
		default:
			from, to, convID = "admin", "insight", "conv-admin-insight"
		}
		seeds = append(seeds, chatSeed{
			ConvID:     convID,
			From:       from,
			To:         to,
			Msg:        fmt.Sprintf("Pesan otomatis #%d — mohon konfirmasi.", i+1),
			MinutesAgo: (i + 1) * 3,
			IsRead:     true,
		})
	}
	return seeds
}

func seedChatMessages(ctx context.Context, repo *repository.MySQLChatRepository, opts Options) error {
	now := time.Now().UTC()
	seeds := buildChatSeeds(opts.ChatMessageCount)
	for i, seed := range seeds {
		_, err := repo.Insert(ctx, model.ChatMessage{
			ConversationID: seed.ConvID,
			FromUser:       seed.From,
			ToUser:         seed.To,
			Message:        seed.Msg,
			Timestamp:      now.Add(-time.Duration(seed.MinutesAgo) * time.Minute),
			IsRead:         seed.IsRead,
		})
		if err != nil {
			return fmt.Errorf("insert chat message #%d: %w", i+1, err)
		}
		if opts.Verbose {
			preview := seed.Msg
			if len(preview) > 40 {
				preview = preview[:40] + "..."
			}
			log.Printf("  [chat_messages] #%d: %s → %s: %q", i+1, seed.From, seed.To, preview)
		}
	}
	return nil
}

// ─── dashboard_data ────────────────────────────────────────────────────────────

func seedDashboardData(ctx context.Context, repo *repository.MySQLDashboardRepository, opts Options) error {
	now := time.Now().UTC()
	cacheTTL := 5 * time.Minute

	caches := []struct {
		key     string
		appName string
		data    any
	}{
		{
			key:     "admin:traffic_summary",
			appName: "",
			data: map[string]any{
				"total_requests":    50,
				"success_count":     42,
				"error_count":       8,
				"success_rate_pct":  84.0,
				"avg_duration_ms":   120,
			},
		},
		{
			key:     "admin:service_indicators",
			appName: "",
			data: map[string]any{
				"Marketplace": map[string]any{"status": "inactive", "last_request": now.Add(-8 * 24 * time.Hour).Format(time.RFC3339)},
				"POS":         map[string]any{"status": "active", "last_request": now.Add(-1 * time.Hour).Format(time.RFC3339)},
				"SupplierHub": map[string]any{"status": "active", "last_request": now.Add(-30 * time.Minute).Format(time.RFC3339)},
				"LogistiKita": map[string]any{"status": "inactive", "last_request": now.Add(-9 * 24 * time.Hour).Format(time.RFC3339)},
				"SmartBank":   map[string]any{"status": "active", "last_request": now.Add(-5 * time.Minute).Format(time.RFC3339)},
			},
		},
		{
			key:     "user:Marketplace:summary",
			appName: "Marketplace",
			data: map[string]any{
				"total_requests":   10,
				"success_count":    8,
				"error_count":      2,
				"success_rate_pct": 80.0,
			},
		},
		{
			key:     "user:POS:summary",
			appName: "POS",
			data: map[string]any{
				"total_requests":   10,
				"success_count":    9,
				"error_count":      1,
				"success_rate_pct": 90.0,
			},
		},
		{
			key:     "monitoring:summary",
			appName: "UMKM Insight",
			data: map[string]any{
				"total_apps":          5,
				"active_apps":         3,
				"inactive_apps":       2,
				"total_requests_week": 50,
				"avg_success_rate":    84.0,
			},
		},
	}

	for i, cache := range caches {
		dataJSON, err := json.Marshal(cache.data)
		if err != nil {
			return fmt.Errorf("marshal dashboard data #%d: %w", i+1, err)
		}
		err = repo.Upsert(ctx, model.DashboardData{
			CacheKey:   cache.key,
			AppName:    cache.appName,
			Data:       dataJSON,
			ComputedAt: now,
			ExpiresAt:  now.Add(cacheTTL),
		})
		if err != nil {
			return fmt.Errorf("upsert dashboard data #%d (%s): %w", i+1, cache.key, err)
		}
		if opts.Verbose {
			log.Printf("  [dashboard_data] #%d: key=%q appName=%q", i+1, cache.key, cache.appName)
		}
	}
	return nil
}

// ─── helpers ───────────────────────────────────────────────────────────────────

func statusText(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "success"
	case code >= 400 && code < 500:
		return "client_error"
	default:
		return "server_error"
	}
}
