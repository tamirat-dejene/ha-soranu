-- +goose Up
CREATE TABLE IF NOT EXISTS user_addresses (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    address VARCHAR(100)
);

-- +goose Down
DROP TABLE IF EXISTS user_addresses;
