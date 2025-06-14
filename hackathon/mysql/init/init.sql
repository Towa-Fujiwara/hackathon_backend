CREATE TABLE users (
                       userId VARCHAR(50) PRIMARY KEY,
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
