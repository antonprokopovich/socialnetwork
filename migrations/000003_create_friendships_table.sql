CREATE TABLE IF NOT EXISTS friendships (
     PRIMARY KEY (user_1_id, user_2_id),
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     user_1_id INT NOT NULL,
     user_2_id INT NOT NULL,
     FOREIGN KEY (user_1_id) REFERENCES users(id) ON DELETE CASCADE,
     FOREIGN KEY (user_2_id) REFERENCES users(id) ON DELETE CASCADE
);