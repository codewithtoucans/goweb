-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    id            SERIAL primary key,
    email         TEXT UNIQUE not null,
    password_hash TEXT        not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
