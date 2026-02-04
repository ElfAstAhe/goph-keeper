package server

import (
	"context"
	"database/sql"
)

// все sql в рамках миграции
const (
	sqlCreateTableUserData string = `
create table if not exists user_data (
    id varchar(50) not null,
    user_id varchar(50) not null,
    name varchar(100) not null,
    kind varchar(50) not null,
    text_data text null,
    binary_data text null,
    created_at timestamptz null default now(),
    modified_at timestamptz null default now(),
    deleted bool not null default false,
    constraint user_data_pk primary key (id),
    constraint user_data_uk unique(user_id, name, kind)
)`
	sqlDropTableUserData string = `drop table if exists user_data cascade`

	sqlCreateIndexUserDataAlive string = `create index user_data_alive_idx on user_data (deleted asc, id asc)`
	sqlDropIndexUserDataAlive   string = `drop index if exists user_data_alive_idx cascade`

	sqlCreateIndexUserDataAliveUser string = `create index user_data_alive_user_idx on user_data (deleted asc, user_id asc)`
	sqlDropIndexUserDataAliveUser   string = `drop index if exists user_data_alive_user_idx cascade`
)

func up00002(ctx context.Context, db *sql.DB) error {

}
