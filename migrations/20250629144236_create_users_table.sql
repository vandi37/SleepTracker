-- +goose Up
-- +goose StatementBegin
create table users (
    id bigserial primary key,
    username varchar(40) unique not null,
    nickname text not null,
    password_hash bytea not null,
    birth date not null,
    created_at timestamptz default current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
