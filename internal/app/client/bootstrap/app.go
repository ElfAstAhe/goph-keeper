// Package bootstrap - application initialization
package bootstrap

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	"github.com/ElfAstAhe/goph-keeper/internal/bll/client/service"
	"github.com/ElfAstAhe/goph-keeper/pkg/client/rest"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/ElfAstAhe/goph-keeper/pkg/logger"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

// App - клиент goph-keeper
type App struct {
	ctx        context.Context
	cancel     context.CancelFunc
	keysHelper *utils.RSAKeysHelper
	settings   *config.AppSettings
	logger     logger.Logger
	client     rest.GophKeeperClient
	cmdHandler service.CmdHandler
	cmdService service.CmdService
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

	// config and params
	app.settings = config.NewAppSettings(app.keysHelper)
	log.Info("loading config")
	if err = app.loadConfig(); err != nil {
		return err
	}

	// logger
	if err = app.initLogger(); err != nil {
		return err
	}

	// command processor
	log.Info("init dependencies")
	if err = app.initDependencies(); err != nil {
		return err
	}

	return nil
}

func (app *App) Run() error {
	if err := app.cmdService.Execute(app.ctx, app.settings.GetCmd(), app.settings.GetOpts()); err != nil {
		return errs.NewAppCommonError("app run", err)
	}

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

func (app *App) initHelpers() error {
	// helpers
	app.keysHelper = utils.NewRSAKeysHelper(utils.RSAKey2048)

	return nil
}

func (app *App) loadConfig() error {
	log := app.logger.GetLogger("bootstrap load config")
	if err := app.settings.Load(); err != nil {
		return errs.NewAppCommonError("load config error", err)
	}

	log.Infof("config FINAL: [%+v]", app.settings)

	// ВНИМАНИЕ! Валидацию пропускаем!

	return nil
}

func (app *App) initLogger() error {
	appLogger, err := logger.NewZapLogger("INFO", "")
	if err != nil {
		return errs.NewAppCommonError("init logger error", err)
	}

	app.logger = appLogger

	return nil
}

func (app *App) initDependencies() error {
	// http client
	app.client = rest.NewGophKeeperSimpleClient(app.settings.GetConfig().Address)

	// services
	app.cmdHandler = service.NewCmdHandlerImpl(app.client, app.settings.GetConfig())
	app.cmdService = service.NewCmdServiceImpl(app.cmdHandler)

	return nil
}
