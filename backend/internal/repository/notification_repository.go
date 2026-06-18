package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

// NotificationRepository mendefinisikan kontrak akses data notifications.
type NotificationRepository interface {
	Insert(ctx context.Context, n model.Notification) (int64, error)
	ListByAppName(ctx context.Context, appName string, limit, offset int) ([]model.Notification, error)
	ListUnread(ctx context.Context, appName string) ([]model.Notification, error)
	MarkAsRead(ctx context.Context, id int64) error
	MarkAllAsRead(ctx context.Context, appName string) error
	CountUnread(ctx context.Context, appName string) (int64, error)
}

// MySQLNotificationRepository mengimplementasikan NotificationRepository dengan MySQL.
type MySQLNotificationRepository struct {
	db *sql.DB
}

// NewMySQLNotificationRepository membuat instance baru.
func NewMySQLNotificationRepository(db *sql.DB) *MySQLNotificationRepository {
	return &MySQLNotificationRepository{db: db}
}

// Insert menyimpan satu notifikasi dan mengembalikan ID-nya.
func (r *MySQLNotificationRepository) Insert(ctx context.Context, n model.Notification) (int64, error) {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO notifications (created_at, app_name, type, message, is_read)
		 VALUES (?, ?, ?, ?, ?)`,
		n.CreatedAt, n.AppName, n.Type, n.Message, boolToInt(n.IsRead),
	)
	if err != nil {
		return 0, fmt.Errorf("insert notification: %w", err)
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// ListByAppName mengambil notifikasi per aplikasi, terbaru dulu.
func (r *MySQLNotificationRepository) ListByAppName(ctx context.Context, appName string, limit, offset int) ([]model.Notification, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, created_at, app_name, type, message, is_read
		 FROM notifications WHERE app_name = ?
		 ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		appName, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list notifications by app name: %w", err)
	}
	defer rows.Close()
	return scanNotifications(rows)
}

// ListUnread mengambil semua notifikasi yang belum dibaca untuk satu aplikasi.
func (r *MySQLNotificationRepository) ListUnread(ctx context.Context, appName string) ([]model.Notification, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, created_at, app_name, type, message, is_read
		 FROM notifications WHERE app_name = ? AND is_read = 0
		 ORDER BY created_at DESC`,
		appName,
	)
	if err != nil {
		return nil, fmt.Errorf("list unread notifications: %w", err)
	}
	defer rows.Close()
	return scanNotifications(rows)
}

// MarkAsRead menandai satu notifikasi sebagai sudah dibaca.
func (r *MySQLNotificationRepository) MarkAsRead(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE notifications SET is_read = 1 WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("mark notification %d as read: %w", id, err)
	}
	return nil
}

// MarkAllAsRead menandai semua notifikasi satu aplikasi sebagai sudah dibaca.
func (r *MySQLNotificationRepository) MarkAllAsRead(ctx context.Context, appName string) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE notifications SET is_read = 1 WHERE app_name = ? AND is_read = 0`,
		appName,
	)
	if err != nil {
		return fmt.Errorf("mark all notifications as read for %s: %w", appName, err)
	}
	return nil
}

// CountUnread menghitung jumlah notifikasi yang belum dibaca untuk satu aplikasi.
func (r *MySQLNotificationRepository) CountUnread(ctx context.Context, appName string) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM notifications WHERE app_name = ? AND is_read = 0`,
		appName,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count unread notifications for %s: %w", appName, err)
	}
	return count, nil
}

// scanNotifications mem-parse baris SQL menjadi slice Notification.
func scanNotifications(rows *sql.Rows) ([]model.Notification, error) {
	var result []model.Notification
	for rows.Next() {
		var n model.Notification
		var isRead int
		if err := rows.Scan(&n.ID, &n.CreatedAt, &n.AppName, &n.Type, &n.Message, &isRead); err != nil {
			return nil, fmt.Errorf("scan notification row: %w", err)
		}
		n.IsRead = isRead == 1
		result = append(result, n)
	}
	return result, rows.Err()
}

// boolToInt mengonversi bool menjadi 0/1 untuk penyimpanan TINYINT MySQL.
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
