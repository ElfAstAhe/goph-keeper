package utils

import (
	"context"
	"database/sql"
)

type DBHelper interface {
	RunInTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error
	ExecStmt(ctx context.Context, stmt *sql.Stmt, notFoundInfo NotFoundInfo, params ...any) error

	IsUniqueViolation(err error) bool
}
