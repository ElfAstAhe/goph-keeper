package server

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

// все sql запросы в рамках этой миграции
const (
	sqlCreateTableUsers string = `create table if not exists users (
    id varchar(50) not null,
    username varchar(50) not null,
    password_hash varchar(2000) null,
    private_key text null,
    public_key text null,
    active bool not null default false,
    deleted bool not null default false,
    person varchar(512) null,
    e_mail varchar(1024) null,
    constraint users_pk primary key (id),
    constraint users_uk unique (username)
)`
	sqlDropTableUsers = `drop table if exists users cascade`

	sqlCreateIndexUsersAlive string = `create index if not exists users_alive_idx on users (deleted asc, id asc)`
	sqlDropIndexUsersAlive   string = `drop index if exists users_alive_idx cascade`

	sqlCreateIndexUsersAliveKey string = `create index if not exists users_alive_key_idx on users (deleted asc, username asc)`
	sqlDropIndexUsersAliveKey   string = `drop index if exists users_alive_key_idx cascade`
)

func up0001(ctx context.Context, db *sql.DB) error {
	if err := upCreateTableUsers(ctx, db); err != nil {
		return err
	}
	if err := upCreateIndexUsersAlive(ctx, db); err != nil {
		return err
	}
	if err := upCreateIndexUsersAliveKey(ctx, db); err != nil {
		return err
	}

	return nil
}

func upCreateTableUsers(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlCreateTableUsers)

	return err
}

func upCreateIndexUsersAlive(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlCreateIndexUsersAlive)

	return err
}

func upCreateIndexUsersAliveKey(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlCreateIndexUsersAliveKey)

	return err
}

func down0001(ctx context.Context, db *sql.DB) error {
	if err := downDropIndexUsersAlive(ctx, db); err != nil {
		return err
	}
	if err := downDropIndexUsersAliveKey(ctx, db); err != nil {
		return err
	}
	if err := downDropTableUsers(ctx, db); err != nil {
		return err
	}

	return nil
}

func downDropTableUsers(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlDropTableUsers)

	return err
}

func downDropIndexUsersAlive(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlDropIndexUsersAlive)

	return err
}

func downDropIndexUsersAliveKey(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlDropIndexUsersAliveKey)

	return err
}

func init() {
	goose.AddMigrationNoTxContext(up0001, down0001)
}
