-- +goose Up
CREATE TABLE notifications (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    app_name   VARCHAR(100)    NOT NULL COMMENT 'Aplikasi yang terkait notifikasi',
    type       ENUM(
        'api_inactive',
        'error_rate',
        'response_time',
        'system'
    )                          NOT NULL DEFAULT 'system',
    message    TEXT            NOT NULL,
    is_read    TINYINT(1)      NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    INDEX idx_notifications_app_name  (app_name),
    INDEX idx_notifications_is_read   (is_read),
    INDEX idx_notifications_created_at(created_at),
    INDEX idx_notifications_type      (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
  COMMENT='Notifikasi sistem untuk admin dan user, termasuk alert API tidak aktif';

-- +goose Down
DROP TABLE notifications;
