package db

import (
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

// Типы поддерживаемых БД
// на данный момент только postgres
const (
	KindInMemory utils.DatabaseKind = "InMemory"
	KindPostgres utils.DatabaseKind = "Postgres"
	KindSQLite3  utils.DatabaseKind = "SQLite3"
)
