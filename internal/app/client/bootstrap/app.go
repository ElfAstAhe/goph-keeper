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

	// helpers
	err := app.initHelpers()
	if err != nil {
		return err
	}

	app.conf = config.NewAppConfig(app.keysHelper)
	log.Info("loading config")
	if err = app.loadConfig(); err != nil {
		return err
	}

	// ToDo: implement

	return nil
}

func (app *App) Run() error {
	// ToDo: implement

	return nil
}

func (app *App) Close() error {
	app.cancel()

	// ToDo: implement

	return nil
}

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

func (app *App) initHelpers() error {
	// helpers
	app.keysHelper = utils.NewRSAKeysHelper(utils.RSAKey2048)

	return nil
}
