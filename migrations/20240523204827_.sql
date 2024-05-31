-- +goose Up
-- +goose StatementBegin
create table tickets(
    id BIGSERIAL PRIMARY KEY NOT NULL ,
    user_email TEXT NOT NULL DEFAULT '',
    price BIGINT NOT NULL,
    home_team TEXT NOT NULL DEFAULT '',
    away_team TEXT NOT NULL DEFAULT '',
    date_time timestamp with time zone default now() not null,
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table tickets;
-- +goose StatementEnd