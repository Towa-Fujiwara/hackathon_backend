SET FOREIGN_KEY_CHECKS = 0;

-- 既存のテーブルを削除します（依存関係の逆順で削除するのが安全です）
-- `users` テーブルが他のテーブルから参照されているため、先に他のテーブルを削除します

DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS follows;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;

-- テーブルを削除した後に、外部キーチェックを再度有効にします
SET FOREIGN_KEY_CHECKS = 1;

CREATE TABLE users (
                       userId VARCHAR(50) PRIMARY KEY,
                       firebaseUid VARCHAR(50) NOT NULL UNIQUE,
                       name VARCHAR(50) NOT NULL UNIQUE,
                       bio TEXT,
                       iconUrl TEXT,
                       createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE posts (
                       id VARCHAR(50) PRIMARY KEY,
                       userId VARCHAR(50) NOT NULL,
                       text TEXT NOT NULL,
                       image TEXT,
                       createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (userId) REFERENCES users(userId)
);

CREATE TABLE likes (
                       id VARCHAR(50) PRIMARY KEY,
                       userId VARCHAR(50) NOT NULL,
                       postId VARCHAR(50) NOT NULL,
                       createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (userId) REFERENCES users(userId),
                       FOREIGN KEY (postId) REFERENCES posts(id)
);

CREATE TABLE comments (
                          id VARCHAR(50) PRIMARY KEY,
                          userId VARCHAR(50) NOT NULL,
                          postId VARCHAR(50) NOT NULL,
                          text TEXT NOT NULL,
                          createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
                          FOREIGN KEY (userId) REFERENCES users(userId),
                          FOREIGN KEY (postId) REFERENCES posts(id)
);

CREATE TABLE follows (
                         id VARCHAR(50) PRIMARY KEY,
                         userId VARCHAR(50) NOT NULL,
                         followUserId VARCHAR(50) NOT NULL,
                         createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
                         FOREIGN KEY (userId) REFERENCES users(userId),
                         FOREIGN KEY (followUserId) REFERENCES users(userId)
);
