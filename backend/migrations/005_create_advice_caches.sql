-- +migrate Up
CREATE TABLE advice_caches (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    cache_date DATE NOT NULL,
    advice TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    UNIQUE KEY uk_user_date (user_id, cache_date),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE advice_caches;
