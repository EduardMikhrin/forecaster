-- +migrate Up

CREATE TABLE subscriptions
(
    email TEXT PRIMARY KEY,
    city_id INTEGER NOT NULL REFERENCES cities(id),
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE cities
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
)

-- +migrate Down

DROP TABLE subscriptions;
DROP TABLE cities;