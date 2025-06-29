-- +goose Up
-- +goose StatementBegin
create type friendship_status as enum ('requested', 'accepted', 'rejected');
create table friends (
    friendship_id bigserial primary key,
    user1_id bigint not null references users(id) on delete cascade,
    user2_id bigint not null references users(id) on delete cascade,
    status friendship_status not null default 'requested',
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    constraint unique_friendship unique (user1_id, user2_id),
    constraint no_self_friendship check (user1_id != user2_id)
);

create index idx_friends_user1_status on friends(user1_id, status);
create index idx_friends_user2_status on friends(user2_id, status); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table friends;
drop type friendship_status;
-- +goose StatementEnd
