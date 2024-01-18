CREATE TABLE IF NOT EXISTS bottles (
       id CHAR(16) NOT NULL PRIMARY KEY,
       msg TEXT NOT NULL,
       token TEXT,
       expiration TEXT NOT NULL,
       expired_at TIMESTAMP,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO bottles (id, msg, token, expiration, expired_at) values
       (1, "", NULL, "300s", NULL),
       (2, "", NULL, "300s", NULL),
       (3, "", NULL, "300s", NULL),
       (4, "", NULL, "300s", NULL),
       (5, "", NULL, "300s", NULL);
