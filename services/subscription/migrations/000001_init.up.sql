CREATE TABLE users(
    id INTEGER PRIMARY KEY,
    user_name TEXT NOT NULL,
    first_name TEXT NOT NULL,
    second_name TEXT NOT NULL
);

CREATE TABLE subscriptions(
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    price INTEGER CHECK ( price > 0 ) NOT NULL,
    currency TEXT CHECK ( currency IN ('RUB','USD','EUR')),
    name TEXT NOT NULL,
    type TEXT CHECK ( type IN ('STREAMING','SOFTWARE','UTILITIES','FINANCE','HEALTH','EDUCATION','OTHER')),
    billing_at TIMESTAMP NOT NULL
)