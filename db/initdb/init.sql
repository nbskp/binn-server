CREATE TABLE bottles (
    id char(255) NOT NULL,
    msg text NOT NULL,
    expired_at timestamp,
    available bool NOT NULL
)
