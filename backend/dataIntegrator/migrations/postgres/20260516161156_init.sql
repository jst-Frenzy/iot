-- +goose Up
-- +goose StatementBegin

create table if not exists devices
(
    name text primary key not null
);

insert into devices (name)
values
    ('temperature_sensor'),
    ('wet_sensor'),
    ('fan'),
    ('pump');

create table if not exists telemetry
(
    id serial primary key,
    device_name text references devices(name),
    value integer not null,
    created_at timestamp default now()
    );

create table if not exists actions_log
(
    id serial primary key,
    device_name text references devices(name),
    action text check (action in ('off', 'on'))
);

create table if not exists user_actions(
    id serial primary key,
    action text check (action in ('off', 'on')),
    device_name text check (device_name in ('fan', 'pump')),
    created_at timestamp default now()
)

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists actions_log;
drop table if exists telemetry;
drop table if exists devices;

-- +goose StatementEnd