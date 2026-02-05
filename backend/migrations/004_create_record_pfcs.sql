-- +migrate Up
CREATE TABLE record_pfcs (
    id VARCHAR(36) PRIMARY KEY,
    record_id VARCHAR(36) NOT NULL UNIQUE,
    protein DOUBLE NOT NULL,
    fat DOUBLE NOT NULL,
    carbs DOUBLE NOT NULL,
    FOREIGN KEY (record_id) REFERENCES records(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE record_pfcs;
