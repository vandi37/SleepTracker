-- +goose Up
-- +goose StatementBegin
create or replace function immutable_date_trunc(text, date)
returns date as $$
begin
    return date_trunc($1, $2);
end;
$$ language plpgsql IMMUTABLE;

create table sleeps (
    id bigserial primary key,
    user_id bigint not null references users(id) on delete cascade,
    sleep_time smallint check (sleep_time >= 0 AND sleep_time < 2160),  
    wake_time smallint check (wake_time >= 0 AND wake_time < 2160),
    score smallint not null check (score >= 0 AND score <= 100),
    enter_date date not null default current_date,
    constraint unique_date unique (user_id, enter_date),
    constraint valid_sleep check (
        (sleep_time is null and wake_time is null) or
        (sleep_time is not null and wake_time is not null and sleep_time < wake_time)
    )
);

create index idx_user_week on sleeps (user_id, immutable_date_trunc('week', enter_date));
create index idx_user_year on sleeps (user_id, immutable_date_trunc('year', enter_date));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table sleeps;
-- +goose StatementEnd
