package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

// ErrCacheNotFound dikembalikan ketika cache key tidak ditemukan.
var ErrCacheNotFound = errors.New("dashboard cache not found")

// DashboardRepository mendefinisikan kontrak akses data dashboard_data.
type DashboardRepository interface {
	Upsert(ctx context.Context, d model.DashboardData) error
	FindByCacheKey(ctx context.Context, cacheKey string) (model.DashboardData, error)
	DeleteExpired(ctx context.Context) (int64, error)
}

// MySQLDashboardRepository mengimplementasikan DashboardRepository dengan MySQL.
type MySQLDashboardRepository struct {
	db *sql.DB
}

// NewMySQLDashboardRepository membuat instance baru.
func NewMySQLDashboardRepository(db *sql.DB) *MySQLDashboardRepository {
	return &MySQLDashboardRepository{db: db}
}

// Upsert menyimpan atau memperbarui cache analytics berdasarkan cache_key.
func (r *MySQLDashboardRepository) Upsert(ctx context.Context, d model.DashboardData) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO dashboard_data (cache_key, app_name, data, computed_at, expires_at)
		 VALUES (?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE
		   app_name    = VALUES(app_name),
		   data        = VALUES(data),
		   computed_at = VALUES(computed_at),
		   expires_at  = VALUES(expires_at)`,
		d.CacheKey, d.AppName, string(d.Data), d.ComputedAt, d.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("upsert dashboard data for key %s: %w", d.CacheKey, err)
	}
	return nil
}

// FindByCacheKey mengambil satu entri cache berdasarkan kunci uniknya.
// Mengembalikan ErrCacheNotFound jika tidak ada.
func (r *MySQLDashboardRepository) FindByCacheKey(ctx context.Context, cacheKey string) (model.DashboardData, error) {
	var d model.DashboardData
	var dataStr string
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, cache_key, app_name, data, computed_at, expires_at
		 FROM dashboard_data WHERE cache_key = ? LIMIT 1`,
		cacheKey,
	).Scan(&d.ID, &d.CacheKey, &d.AppName, &dataStr, &d.ComputedAt, &d.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return model.DashboardData{}, ErrCacheNotFound
	}
	if err != nil {
		return model.DashboardData{}, fmt.Errorf("find dashboard data by key %s: %w", cacheKey, err)
	}
	d.Data = []byte(dataStr)
	return d, nil
}

// DeleteExpired menghapus semua entri cache yang sudah melewati expires_at.
// Mengembalikan jumlah baris yang dihapus.
func (r *MySQLDashboardRepository) DeleteExpired(ctx context.Context) (int64, error) {
	result, err := r.db.ExecContext(
		ctx,
		`DELETE FROM dashboard_data WHERE expires_at < ?`,
		time.Now().UTC(),
	)
	if err != nil {
		return 0, fmt.Errorf("delete expired dashboard cache: %w", err)
	}
	n, _ := result.RowsAffected()
	return n, nil
}
