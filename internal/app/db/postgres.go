package db

import (
	"database/sql"
	"time"

	"github.com/ElfAstAhe/goph-keeper/internal/app/config"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresDB struct {
	db   *sql.DB
	kind utils.DatabaseKind
	dsn  string
}

// NewPostgresDB - конструктор соединения с БД postgres
func NewPostgresDB(dbConf *config.DatabaseConfig) (*PostgresDB, error) {
	err := utils.DBValidateDSN(dbConf.DSN)
	if err != nil {
		return nil, err
	}

	pg, err := sql.Open("pgx", dbConf.DSN)
	if err != nil {
		return nil, err
	}

	pg.SetMaxOpenConns(dbConf.MaxOpenConnections)
	pg.SetMaxIdleConns(dbConf.MaxIdleConnections)
	pg.SetConnMaxIdleTime(time.Duration(dbConf.MaxIdleConnectionLifetime) * time.Second)

	err = pg.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresDB{
		db:   pg,
		kind: KindPostgres,
		dsn:  dbConf.DSN,
	}, nil
}

// io.Closer =====================

func (db *PostgresDB) Close() error {
	return db.db.Close()
}

// db.DB =========================

func (db *PostgresDB) GetDBKind() utils.DatabaseKind {
	return db.kind
}

func (db *PostgresDB) GetDsn() string {
	return db.dsn
}

func (db *PostgresDB) GetDB() *sql.DB {
	return db.db
}
