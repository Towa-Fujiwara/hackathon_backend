CREATE TABLE users (
                       id VARCHAR(50) PRIMARY KEY,
                       name VARCHAR(50) NOT NULL UNIQUE,
                       password TEXT NOT NULL,
                       display_name TEXT NOT NULL,
                       bio TEXT,
                       icon_url TEXT,
                       age INTEGER,
                       created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE posts (
                       id VARCHAR(50) PRIMARY KEY,
                       user_id VARCHAR(50) NOT NULL,
                       text TEXT NOT NULL,
                       image TEXT,
                       created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE likes (
                       id VARCHAR(50) PRIMARY KEY,
                       user_id VARCHAR(50) NOT NULL,
                       post_id VARCHAR(50) NOT NULL,
                       created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (user_id) REFERENCES users(id),
                       FOREIGN KEY (post_id) REFERENCES posts(id)
);

CREATE TABLE comments (
                          id VARCHAR(50) PRIMARY KEY,
                          user_id VARCHAR(50) NOT NULL,
                          post_id VARCHAR(50) NOT NULL,
                          text TEXT NOT NULL,
                          created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                          FOREIGN KEY (user_id) REFERENCES users(id),
                          FOREIGN KEY (post_id) REFERENCES posts(id)
);

CREATE TABLE follows (
                         id VARCHAR(50) PRIMARY KEY,
                         user_id VARCHAR(50) NOT NULL,
                         follow_user_id VARCHAR(50) NOT NULL,
                         created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                         FOREIGN KEY (user_id) REFERENCES users(id),
                         FOREIGN KEY (follow_user_id) REFERENCES users(id)
);
