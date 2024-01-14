-- +goose Up
-- +goose StatementBegin
create table sessions(
    id         serial primary key,
    user_id    int unique,
    token_hash text not null,
    foreign key (user_id) references users (id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
