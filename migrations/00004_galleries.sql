-- +goose Up
-- +goose StatementBegin
CREATE TABLE galleries (
  id SERIAL PRIMARY KEY,
  user_id int REFERENCES users(id),
  title VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE galleries;
-- +goose StatementEnd
