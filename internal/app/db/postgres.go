package db

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresDB struct {
	db   *sql.DB
	kind Kind
	dsn  string
}

// NewPostgresDB - конструктор соединения с БД postgres
func NewPostgresDB(dsn string) (*PostgresDB, error) {
	pg, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	pg.SetMaxOpenConns(20)
	pg.SetMaxIdleConns(5)
	pg.SetConnMaxIdleTime(60 * time.Second)

	err = pg.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresDB{
		db:   pg,
		kind: KindPostgres,
		dsn:  dsn,
	}, nil
}

// io.Closer =====================

func (db *PostgresDB) Close() error {
	return db.db.Close()
}

// db.DB =========================

func (db *PostgresDB) Kind() Kind {
	return db.kind
}

func (db *PostgresDB) Dsn() string {
	return db.dsn
}

func (db *PostgresDB) GetDB() *sql.DB {
	return db.db
}
