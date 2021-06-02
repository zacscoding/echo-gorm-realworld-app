CREATE TABLE users
(
    user_id    int unsigned auto_increment PRIMARY KEY,
    email      VARCHAR(255) NOT NULL,
    name       VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    bio        TEXT,
    image      VARCHAR(255) NULL,
    created_at datetime NULL,
    updated_at datetime NULL,
    deleted_at datetime NULL,
    disabled   TINYINT(1) default 0,
    UNIQUE KEY unique_users_email (email)
) CHARACTER SET utf8mb4;
CREATE INDEX idx_users_name ON users (name);

CREATE TABLE follows
(
    created_at DATETIME NULL,
    user_id    INT UNSIGNED,
    follow_id  INT UNSIGNED,
    CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users (user_id),
    CONSTRAINT follow_id_fk FOREIGN KEY (follow_id) REFERENCES users (user_id),
    PRIMARY KEY (user_id, follow_id)
) CHARACTER SET utf8mb4;