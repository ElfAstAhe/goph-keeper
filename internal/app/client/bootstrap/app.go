// Package bootstrap - application initialization
package bootstrap

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/ElfAstAhe/goph-keeper/pkg/logger"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

// App - клиент goph-keeper
type App struct {
	ctx        context.Context
	cancel     context.CancelFunc
	keysHelper *utils.RSAKeysHelper
	conf       *config.AppConfig
	logger     logger.Logger
}

func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		ctx:    ctx,
		cancel: cancel,
		logger: logger.NewStartupZapLogger(),
	}
}

func (app *App) Init() error {
	log := app.logger.GetLogger("bootstrap init")

	app.conf = config.NewConfig()
	log.Info("loading config")
	if err := app.loadConfig(); err != nil {
		return err
	}
}

func (app *App) Run() error {

}

func (app *App) Close() error {}

func (app *App) GetLogger() logger.Logger {
	return app.logger
}

func (app *App) loadConfig() error {
	log := app.logger.GetLogger("bootstrap load config")
	if err := app.conf.Load(); err != nil {
		return errs.NewAppCommonError("load config error", err)
	}

	log.Infof("config FINAL: [%+v]", app.conf)

	if err := app.conf.Validate(); err != nil {
		return errs.NewAppConfigError("validate config error", err)
	}

	return nil
}
