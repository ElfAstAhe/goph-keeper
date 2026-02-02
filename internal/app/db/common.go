package db

import (
	"database/sql"
	"io"
)

type Kind string

const (
	KindInMemory Kind = "InMemory"
	KindPostgres Kind = "Postgres"
)

type DB interface {
	GetDB() *sql.DB
	GetDBKind() string
	GetDsn() string
}

func CloseDB(db DB) error {
	if closer, ok := db.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
