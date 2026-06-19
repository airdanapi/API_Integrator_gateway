package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

// LogRepository mendefinisikan kontrak akses data request_logs.
type LogRepository interface {
	Insert(ctx context.Context, log model.RequestLog) (int64, error)
	ListBySourceApp(ctx context.Context, sourceApp string, limit, offset int) ([]model.RequestLog, error)
	ListRecent(ctx context.Context, limit, offset int) ([]model.RequestLog, error)
	CountByStatus(ctx context.Context, since time.Time) (map[int]int64, error)
	CountByStatusForApp(ctx context.Context, appName string, since time.Time) (map[int]int64, error)
	CountBySourceApp(ctx context.Context, since time.Time) (map[string]int64, error)
}

// MySQLLogRepository mengimplementasikan LogRepository menggunakan MySQL.
type MySQLLogRepository struct {
	db *sql.DB
}

// NewMySQLLogRepository membuat instance baru MySQLLogRepository.
func NewMySQLLogRepository(db *sql.DB) *MySQLLogRepository {
	return &MySQLLogRepository{db: db}
}

// Insert menyimpan satu request log dan mengembalikan ID yang baru dibuat.
func (r *MySQLLogRepository) Insert(ctx context.Context, log model.RequestLog) (int64, error) {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO request_logs (timestamp, source_app, endpoint, method, payload, status, response, duration_ms)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		log.Timestamp,
		log.SourceApp,
		log.Endpoint,
		log.Method,
		nullableJSON(log.Payload),
		log.Status,
		nullableJSON(log.Response),
		log.DurationMS,
	)
	if err != nil {
		return 0, fmt.Errorf("insert request log: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id: %w", err)
	}
	return id, nil
}

// ListBySourceApp mengambil log berdasarkan aplikasi sumber dengan pagination.
func (r *MySQLLogRepository) ListBySourceApp(ctx context.Context, sourceApp string, limit, offset int) ([]model.RequestLog, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, timestamp, source_app, endpoint, method, payload, status, response, duration_ms
		 FROM request_logs WHERE source_app = ?
		 ORDER BY timestamp DESC LIMIT ? OFFSET ?`,
		sourceApp, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list request logs by source app: %w", err)
	}
	defer rows.Close()
	return scanRequestLogs(rows)
}

// ListRecent mengambil log terbaru (semua aplikasi) dengan pagination.
func (r *MySQLLogRepository) ListRecent(ctx context.Context, limit, offset int) ([]model.RequestLog, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, timestamp, source_app, endpoint, method, payload, status, response, duration_ms
		 FROM request_logs ORDER BY timestamp DESC LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list recent request logs: %w", err)
	}
	defer rows.Close()
	return scanRequestLogs(rows)
}

// CountByStatus menghitung jumlah log per HTTP status code sejak waktu tertentu.
func (r *MySQLLogRepository) CountByStatus(ctx context.Context, since time.Time) (map[int]int64, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT status, COUNT(*) FROM request_logs WHERE timestamp >= ? GROUP BY status`,
		since,
	)
	if err != nil {
		return nil, fmt.Errorf("count request logs by status: %w", err)
	}
	defer rows.Close()
	result := make(map[int]int64)
	for rows.Next() {
		var status int
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("scan status count row: %w", err)
		}
		result[status] = count
	}
	return result, rows.Err()
}

// CountByStatusForApp menghitung jumlah log per HTTP status code untuk satu aplikasi sejak waktu tertentu.
func (r *MySQLLogRepository) CountByStatusForApp(ctx context.Context, appName string, since time.Time) (map[int]int64, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT status, COUNT(*) FROM request_logs WHERE source_app = ? AND timestamp >= ? GROUP BY status`,
		appName, since,
	)
	if err != nil {
		return nil, fmt.Errorf("count request logs by status for app: %w", err)
	}
	defer rows.Close()
	result := make(map[int]int64)
	for rows.Next() {
		var status int
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("scan status for app count row: %w", err)
		}
		result[status] = count
	}
	return result, rows.Err()
}

// CountBySourceApp menghitung jumlah log per aplikasi sumber sejak waktu tertentu.
func (r *MySQLLogRepository) CountBySourceApp(ctx context.Context, since time.Time) (map[string]int64, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT source_app, COUNT(*) FROM request_logs WHERE timestamp >= ? GROUP BY source_app`,
		since,
	)
	if err != nil {
		return nil, fmt.Errorf("count request logs by source app: %w", err)
	}
	defer rows.Close()
	result := make(map[string]int64)
	for rows.Next() {
		var app string
		var count int64
		if err := rows.Scan(&app, &count); err != nil {
			return nil, fmt.Errorf("scan source app count row: %w", err)
		}
		result[app] = count
	}
	return result, rows.Err()
}

// scanRequestLogs mem-parse baris SQL menjadi slice RequestLog.
func scanRequestLogs(rows *sql.Rows) ([]model.RequestLog, error) {
	var logs []model.RequestLog
	for rows.Next() {
		var l model.RequestLog
		var payload, response sql.NullString
		var durationMS sql.NullInt32
		if err := rows.Scan(
			&l.ID, &l.Timestamp, &l.SourceApp, &l.Endpoint, &l.Method,
			&payload, &l.Status, &response, &durationMS,
		); err != nil {
			return nil, fmt.Errorf("scan request log row: %w", err)
		}
		if payload.Valid {
			l.Payload = []byte(payload.String)
		}
		if response.Valid {
			l.Response = []byte(response.String)
		}
		if durationMS.Valid {
			v := int(durationMS.Int32)
			l.DurationMS = &v
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

// nullableJSON mengonversi []byte menjadi sql.NullString.
// Digunakan untuk kolom JSON nullable di database.
func nullableJSON(b []byte) sql.NullString {
	if len(b) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{String: string(b), Valid: true}
}
