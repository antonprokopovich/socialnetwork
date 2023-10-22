ALTER table users ADD COLUMN email VARCHAR(255) NOT NULL;
ALTER table users ADD COLUMN hashed_password CHAR(60) NOT NULL;

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);