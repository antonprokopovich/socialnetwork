ALTER TABLE users
MODIFY first_name VARCHAR(255) NOT NULL,
MODIFY last_name VARCHAR(255) NOT NULL;

CREATE INDEX idx_first_name ON users(first_name) USING BTREE;
CREATE INDEX idx_last_name ON users(last_name) USING BTREE;