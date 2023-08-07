CREATE TABLE bottles (
    id char(255) NOT NULL PRIMARY KEY,
    msg text NOT NULL,
    expired_at timestamp,
    available bool NOT NULL
);

CREATE TABLE subscriptions (
    id char(255) NOT NULL PRIMARY KEY,
    expired_at timestamp NOT NULL,
    next_time timestamp NOT NULL
);

CREATE TABLE subscribed_bottles (
       bottle_id char(255) NOT NULL PRIMARY KEY,
       subscription_id char(255) NOT NULL,
       FOREIGN KEY fk_bottle_id (bottle_id)
       REFERENCES bottles(id),
       FOREIGN KEY fk_subscription_id (subscription_id)
       REFERENCES subscriptions(id)
);

CREATE TABLE tokens (
       id char(255) NOT NULL PRIMARY KEY
)
