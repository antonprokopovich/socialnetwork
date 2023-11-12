CREATE TABLE IF NOT EXISTS friend_requests (
    PRIMARY KEY (sender_user_id, recipient_user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sender_user_id INT NOT NULL,
    recipient_user_id INT NOT NULL,
    FOREIGN KEY (sender_user_id) REFERENCES users(id),
    FOREIGN KEY (recipient_user_id) REFERENCES users(id)
);