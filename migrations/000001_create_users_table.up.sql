CREATE TABLE IF NOT EXISTS users (
    id int PRIMARY KEY AUTO_INCREMENT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    first_name text NOT NULL,
    last_name text NOT NULL,
    age integer DEFAULT 0 NOT NULL,
    gender text NOT NULL,
    interests text,
    city text NOT NULL
);
