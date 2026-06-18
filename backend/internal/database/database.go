package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/config"
	"github.com/go-sql-driver/mysql"
)

func Open(ctx context.Context, cfg config.Config) (*sql.DB, error) {
	mysqlConfig := mysql.Config{
		User:                 cfg.DBUser,
		Passwd:               cfg.DBPassword,
		Net:                  "tcp",
		Addr:                 cfg.DBHost + ":" + cfg.DBPort,
		DBName:               cfg.DBName,
		ParseTime:            true,
		Loc:                  time.UTC,
		AllowNativePasswords: true,
	}
	dsn := mysqlConfig.FormatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open MySQL connection: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(time.Minute)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping MySQL: %w", err)
	}
	return db, nil
}
