-- +goose Up
CREATE TABLE request_logs (
    id            BIGINT UNSIGNED   NOT NULL AUTO_INCREMENT,
    timestamp     DATETIME(3)       NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    source_app    VARCHAR(100)      NOT NULL COMMENT 'Nama aplikasi pemanggil, misal Marketplace',
    endpoint      VARCHAR(255)      NOT NULL COMMENT 'Endpoint tujuan yang di-forward',
    method        VARCHAR(10)       NOT NULL DEFAULT 'POST',
    payload       JSON              NULL     COMMENT 'Request payload JSON yang dikirim',
    status        SMALLINT UNSIGNED NOT NULL COMMENT 'HTTP status code respon dari tujuan',
    response      JSON              NULL     COMMENT 'Response body JSON dari tujuan',
    duration_ms   INT UNSIGNED      NULL     COMMENT 'Durasi request dalam milidetik',
    PRIMARY KEY (id),
    INDEX idx_request_logs_timestamp    (timestamp),
    INDEX idx_request_logs_source_app   (source_app),
    INDEX idx_request_logs_endpoint     (endpoint),
    INDEX idx_request_logs_status       (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
  COMMENT='Audit log setiap request yang masuk dan diproses gateway';

-- +goose Down
DROP TABLE request_logs;
