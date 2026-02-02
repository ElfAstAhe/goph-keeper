package bootstrap

import (
	"context"
	"goph-keeper/internal/app/db"
)

// App - приложение
type App struct {
	ctx    context.Context
	Cancel context.CancelFunc
	db     db.DB
}
