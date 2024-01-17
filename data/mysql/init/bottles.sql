CREATE TABLE IF NOT EXISTS bottles (
       id CHAR(16) NOT NULL PRIMARY KEY,
       msg TEXT NOT NULL,
       expired_at TIMESTAMP,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO bottles (id, msg, expired_at) values
       (1, "", NULL),
       (2, "", NULL),
       (3, "", NULL),
       (4, "", NULL),
       (5, "", NULL),
       (6, "", NULL),
       (7, "", NULL),
       (8, "", NULL),
       (9, "", NULL),
       (10, "", NULL);
