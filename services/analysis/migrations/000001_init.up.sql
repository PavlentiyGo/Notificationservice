
CREATE TYPE subscription_type AS ENUM ('STREAMING','SOFTWARE','UTILITIES','FINANCE','HEALTH','EDUCATION','OTHER');

CREATE TYPE currency AS ENUM ('RUB','USD','EUR');

CREATE TABLE payments(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    date TIMESTAMP WITH TIME ZONE,
    subscription_name TEXT NOT NULL,
    subscription_type subscription_type,
    subscription_currency currency,
    price DECIMAL
);