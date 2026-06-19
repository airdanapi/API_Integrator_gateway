package dashboard

import (
	"context"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

const inactiveThreshold = 7 * 24 * time.Hour

// knownApps adalah daftar aplikasi yang dipantau gateway sesuai PRD.
var knownApps = []string{
	"Marketplace", "POS", "SupplierHub", "LogistiKita", "SmartBank",
}

// TrafficSummary merangkum statistik traffic gateway dalam periode tertentu.
type TrafficSummary struct {
	TotalRequests  int64   `json:"total_requests"`
	SuccessCount   int64   `json:"success_count"`
	ErrorCount     int64   `json:"error_count"`
	SuccessRatePct float64 `json:"success_rate_pct"`
	AvgDurationMS  int     `json:"avg_duration_ms"`
}

// ServiceIndicator menunjukkan status aktif/inaktif setiap aplikasi.
type ServiceIndicator struct {
	AppName     string    `json:"app_name"`
	Status      string    `json:"status"` // "active" | "inactive"
	LastRequest time.Time `json:"last_request"`
}

// AuditLogEntry adalah satu entri log untuk ditampilkan di tabel audit.
type AuditLogEntry struct {
	ID         int64     `json:"id"`
	Timestamp  time.Time `json:"timestamp"`
	SourceApp  string    `json:"source_app"`
	Endpoint   string    `json:"endpoint"`
	Method     string    `json:"method"`
	Status     int       `json:"status"`
	DurationMS *int      `json:"duration_ms"`
}

// LogQuerier mendefinisikan query yang diperlukan dari repository log.
type LogQuerier interface {
	ListRecent(ctx context.Context, limit, offset int) ([]model.RequestLog, error)
	CountByStatus(ctx context.Context, since time.Time) (map[int]int64, error)
	CountBySourceApp(ctx context.Context, since time.Time) (map[string]int64, error)
}

// Service menyediakan logika bisnis dashboard admin.
type Service struct {
	logs LogQuerier
}

// New membuat instance Service baru.
func New(logs LogQuerier) *Service {
	return &Service{logs: logs}
}

// GetTrafficSummary menghitung ringkasan traffic sejak waktu `since`.
func (s *Service) GetTrafficSummary(ctx context.Context, since time.Time) (TrafficSummary, error) {
	counts, err := s.logs.CountByStatus(ctx, since)
	if err != nil {
		return TrafficSummary{}, err
	}

	var total, success, errCount int64
	for status, count := range counts {
		total += count
		if status >= 200 && status < 300 {
			success += count
		} else {
			errCount += count
		}
	}

	var ratePct float64
	if total > 0 {
		ratePct = float64(success) / float64(total) * 100
	}

	return TrafficSummary{
		TotalRequests:  total,
		SuccessCount:   success,
		ErrorCount:     errCount,
		SuccessRatePct: ratePct,
		AvgDurationMS:  0, // placeholder; avg query bisa ditambah di Sprint berikutnya
	}, nil
}

// GetServiceIndicators mendeteksi aplikasi yang aktif atau inaktif (>= 1 minggu tanpa request).
func (s *Service) GetServiceIndicators(ctx context.Context) ([]ServiceIndicator, error) {
	since := time.Now().UTC().Add(-inactiveThreshold)
	countsByApp, err := s.logs.CountBySourceApp(ctx, since)
	if err != nil {
		return nil, err
	}

	// Untuk lastRequest, kita gunakan now jika aktif, atau now-threshold jika inaktif (simulasi).
	// Query last_request per app bisa ditambahkan di Sprint berikutnya jika dibutuhkan.
	now := time.Now().UTC()
	indicators := make([]ServiceIndicator, 0, len(knownApps))
	for _, app := range knownApps {
		status := "inactive"
		lastReq := now.Add(-inactiveThreshold - time.Hour) // default: lebih dari threshold
		if count, ok := countsByApp[app]; ok && count > 0 {
			status = "active"
			lastReq = now.Add(-time.Hour) // aktif: simulasi 1 jam lalu
		}
		indicators = append(indicators, ServiceIndicator{
			AppName:     app,
			Status:      status,
			LastRequest: lastReq,
		})
	}
	return indicators, nil
}

// GetAuditLogs mengambil log terbaru dengan pagination.
// Mengembalikan (items, total, error).
func (s *Service) GetAuditLogs(ctx context.Context, limit, offset int) ([]AuditLogEntry, int64, error) {
	logs, err := s.logs.ListRecent(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Hitung total dari CountByStatus (semua waktu) sebagai estimasi total log
	counts, err := s.logs.CountByStatus(ctx, time.Time{}) // zero time = semua log
	if err != nil {
		return nil, 0, err
	}
	var total int64
	for _, c := range counts {
		total += c
	}

	entries := make([]AuditLogEntry, 0, len(logs))
	for _, l := range logs {
		entries = append(entries, AuditLogEntry{
			ID:         l.ID,
			Timestamp:  l.Timestamp,
			SourceApp:  l.SourceApp,
			Endpoint:   l.Endpoint,
			Method:     l.Method,
			Status:     l.Status,
			DurationMS: l.DurationMS,
		})
	}
	return entries, total, nil
}

