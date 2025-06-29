-- +goose Up
-- +goose StatementBegin
create table sleeps (
    id bigserial primary key,
    user_id bigint not null references users(id),
    sleep_time smallint check (sleep_time >= 0 AND sleep_time < 1440),  
    wake_time smallint check (wake_time >= 0 AND wake_time < 1440),
    score smallint not null check (score >= 0 AND score <= 100),
    enter_date date not null default current_date,
    constraint unique_date unique (user_id, enter_date),
    constraint valid_sleep check (
        (sleep_time is null and wake_time is null) or
        (sleep_time is not null and wake_time is not null)
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table sleeps;
-- +goose StatementEnd
