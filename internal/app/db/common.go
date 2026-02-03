package db

import (
	"database/sql"
	"io"
	"strings"
)

type Kind string

const (
	KindInMemory Kind = "InMemory"
	KindPostgres Kind = "Postgres"
	KindSQLite3  Kind = "SQLite3"
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

func DBKindFromDSN(dsn string) Kind {
	switch {
	case strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://"):
		return KindPostgres
	case strings.Contains(dsn, ".db") || strings.Contains(dsn, "mode=memory"):
		return KindSQLite3
	default:
		return KindInMemory
	}
}
