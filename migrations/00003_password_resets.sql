-- +goose Up
-- +goose StatementBegin
create table password_resets(
    id         serial primary key,
    user_id    int unique,
    token_hash text not null,
    expires_at timestamptz not null,
    foreign key (user_id) references users (id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table password_resets;
-- +goose StatementEnd
