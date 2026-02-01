-- +migrate Up
CREATE TABLE records (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    eaten_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_records_user_id (user_id),
    INDEX idx_records_eaten_at (eaten_at),
    INDEX idx_records_user_eaten_at (user_id, eaten_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE record_items (
    id VARCHAR(36) PRIMARY KEY,
    record_id VARCHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    calories INT NOT NULL,
    INDEX idx_record_items_record_id (record_id),
    FOREIGN KEY (record_id) REFERENCES records(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE record_items;
DROP TABLE records;
