package server

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
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
    constraint user_data_uk unique(user_id, name, kind),
    constraint user_data_user_fk foreign key (user_id) references users(id)
)`
	sqlDropTableUserData string = `drop table if exists user_data cascade`

	sqlCreateIndexUserDataAlive string = `create index user_data_alive_idx on user_data (deleted asc, id asc)`
	sqlDropIndexUserDataAlive   string = `drop index if exists user_data_alive_idx cascade`

	sqlCreateIndexUserDataAliveUser string = `create index user_data_alive_user_idx on user_data (deleted asc, user_id asc)`
	sqlDropIndexUserDataAliveUser   string = `drop index if exists user_data_alive_user_idx cascade`
)

func up00002(ctx context.Context, db *sql.DB) error {
	if err := createTableUserData(ctx, db); err != nil {
		return err
	}
	if err := createIndexUserDataAlive(ctx, db); err != nil {
		return err
	}
	if err := createIndexUserDataAliveUser(ctx, db); err != nil {
		return err
	}

	return nil
}

func createTableUserData(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlCreateTableUserData)

	return err
}

func createIndexUserDataAlive(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlCreateIndexUserDataAlive)

	return err
}

func createIndexUserDataAliveUser(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlCreateIndexUserDataAliveUser)

	return err
}

func down00002(ctx context.Context, db *sql.DB) error {
	if err := dropIndexUserDataAlive(ctx, db); err != nil {
		return err
	}
	if err := dropIndexUserDataAliveUser(ctx, db); err != nil {
		return err
	}
	if err := dropTableUserData(ctx, db); err != nil {
		return err
	}

	return nil
}

func dropTableUserData(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlDropTableUserData)

	return err
}

func dropIndexUserDataAlive(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlDropIndexUserDataAlive)

	return err
}

func dropIndexUserDataAliveUser(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlDropIndexUserDataAliveUser)

	return err
}

func init() {
	goose.AddMigrationNoTxContext(up00002, down00002)
}
