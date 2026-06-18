-- +goose Up
CREATE TABLE dashboard_data (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    cache_key    VARCHAR(255)    NOT NULL COMMENT 'Kunci unik untuk data yang di-cache',
    app_name     VARCHAR(100)    NOT NULL DEFAULT '' COMMENT 'Kosong untuk data global, diisi untuk per-aplikasi',
    data         JSON            NOT NULL COMMENT 'Payload analytics/cache JSON',
    computed_at  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT 'Waktu komputasi data ini',
    expires_at   DATETIME(3)     NOT NULL COMMENT 'Waktu kedaluwarsa cache ini',
    PRIMARY KEY (id),
    UNIQUE INDEX uq_dashboard_data_cache_key (cache_key),
    INDEX idx_dashboard_data_app_name   (app_name),
    INDEX idx_dashboard_data_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
  COMMENT='Cache analytics dan data agregat untuk dashboard admin dan user';

-- +goose Down
DROP TABLE dashboard_data;
