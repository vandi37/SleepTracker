-- +goose Up
-- +goose StatementBegin
create table friends (
    id bigserial primary key,
    user1_id bigint not null references users(id) on delete cascade,
    user2_id bigint not null references users(id) on delete cascade,
    is_accepted boolean not null default false,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp,
    constraint unique_friendship unique (user1_id, user2_id),
    constraint no_self_friendship check (user1_id != user2_id)
);

create index idx_friends_user1_status on friends(user1_id, is_accepted);
create index idx_friends_user2_status on friends(user2_id, is_accepted); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table friends;
-- +goose StatementEnd
