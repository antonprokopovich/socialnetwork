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

# Test data

INSERT INTO users (id, created_at, first_name, last_name, age, gender, interests, city) VALUES (
    0,
    CURRENT_TIMESTAMP,
    'John',
    'Doe',
    18,
    'male',
    'sports, books, science',
    'London'
);