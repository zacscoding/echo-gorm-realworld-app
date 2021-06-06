-- -----------------------------------------------------
-- users
-- -----------------------------------------------------
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

-- -----------------------------------------------------
-- follow
-- -----------------------------------------------------
CREATE TABLE follows
(
    created_at DATETIME NULL,
    user_id    INT UNSIGNED,
    follow_id  INT UNSIGNED,
    CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users (user_id),
    CONSTRAINT follow_id_fk FOREIGN KEY (follow_id) REFERENCES users (user_id),
    PRIMARY KEY (user_id, follow_id)
) CHARACTER SET utf8mb4;

-- -----------------------------------------------------
-- articles
-- -----------------------------------------------------
-- TODO: add indices
CREATE TABLE articles
(
    article_id  INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at  DATETIME NULL,
    updated_at  DATETIME NULL,
    deleted_at  DATETIME NULL,
    slug        VARCHAR(255),
    title       VARCHAR(255),
    description TEXT,
    body        TEXT,
    author_id   INT UNSIGNED,
    UNIQUE KEY unique_articles_slug (slug),
    CONSTRAINT articles_author_id_fk
        FOREIGN KEY (author_id) REFERENCES users (user_id)
) CHARACTER SET utf8mb4;
CREATE INDEX idx_articles_deleted_at ON articles (deleted_at);

-- -----------------------------------------------------
-- article_favorites
-- -----------------------------------------------------
CREATE TABLE article_favorites
(
    user_id    INT UNSIGNED,
    article_id INT UNSIGNED,
    CONSTRAINT article_favorites_user_id_fk
        FOREIGN KEY (user_id) REFERENCES users (user_id),
    CONSTRAINT article_favorites_article_id_fk
        FOREIGN KEY (article_id) REFERENCES articles (article_id),
    PRIMARY KEY (user_id, article_id)
);

-- -----------------------------------------------------
-- tags
-- -----------------------------------------------------
CREATE TABLE tags
(
    tag_id         INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NULL,
    name       varchar(255),
    UNIQUE KEY unique_tags_name (name)
) CHARACTER SET utf8mb4;

-- -----------------------------------------------------
-- tags
-- -----------------------------------------------------
CREATE TABLE article_tags
(
    article_id INT UNSIGNED,
    tag_id     INT UNSIGNED,
    CONSTRAINT article_tags_article_id_fk
        FOREIGN KEY (article_id) REFERENCES articles (article_id),
    CONSTRAINT article_tags_tag_id_fk
        FOREIGN KEY (tag_id) REFERENCES tags (tag_id),
    PRIMARY KEY (article_id, tag_id)
);

-- -----------------------------------------------------
-- comments
-- -----------------------------------------------------
CREATE TABLE comments
(
    comment_id         INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NULL,
    updated_at DATETIME NULL,
    deleted_at datetime NULL,
    body       text,
    article_id INT UNSIGNED,
    author_id  INT UNSIGNED,
    CONSTRAINT comments_article_id_fk
        FOREIGN KEY (article_id) REFERENCES articles (article_id),
    CONSTRAINT comments_author_id_fk
        FOREIGN KEY (author_id) REFERENCES users (user_id)
) CHARACTER SET utf8mb4;
CREATE index idx_comments_deleted_at ON comments (deleted_at);