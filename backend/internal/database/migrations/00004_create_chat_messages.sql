-- +goose Up
CREATE TABLE chat_messages (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    conversation_id VARCHAR(100)    NOT NULL COMMENT 'ID percakapan untuk mengelompokkan pesan',
    from_user       VARCHAR(100)    NOT NULL COMMENT 'Username pengirim',
    to_user         VARCHAR(100)    NOT NULL COMMENT 'Username penerima',
    message         TEXT            NOT NULL,
    timestamp       DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    is_read         TINYINT(1)      NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    INDEX idx_chat_messages_conversation_id (conversation_id),
    INDEX idx_chat_messages_from_user       (from_user),
    INDEX idx_chat_messages_to_user         (to_user),
    INDEX idx_chat_messages_timestamp       (timestamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
  COMMENT='Pesan chat antara admin gateway dan user aplikasi';

-- +goose Down
DROP TABLE chat_messages;
