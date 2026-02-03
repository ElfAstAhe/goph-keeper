package utils

import (
	"database/sql"
	"io"

	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/xo/dburl"
)

type DatabaseKind string

type DB interface {
	GetDB() *sql.DB
	GetDBKind() DatabaseKind
	GetDsn() string
}

func DBClose(db DB) error {
	if closer, ok := db.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}

func DBValidateDSN(dsn string) error {
	_, err := dburl.Parse(dsn)
	if err != nil {
		return errs.NewUtlDBError("dsn parse error", err)
	}

	return nil
}
