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
	ctx          context.Context
	cancel       context.CancelFunc
	keysHelper   *utils.RSAKeysHelper
	conf         *config.AppConfig
	logger       logger.Logger
	client       rest.GophKeeperClient
	cmdProcessor *service.CommandProcessor
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
	app.conf = config.NewAppConfig(app.keysHelper)
	log.Info("loading config")
	if err = app.loadConfig(); err != nil {
		return err
	}

	// logger
	if err = app.initLogger(); err != nil {
		return err
	}

	// http client
	log.Info("http client")
	app.client = rest.NewGophKeeperSimpleClient(app.conf.Address)

	// command processor
	log.Info("command processor")
	if err = app.initCommandProcessor(); err != nil {
		return err
	}

	return nil
}

func (app *App) initLogger() error {
	// ToDo: implement

	return nil
}

func (app *App) initCommandProcessor() error {
	app.cmdProcessor = service.NewCommandProcessor(map[service.Command]service.ProcessorFunc{
		service.CmdGenConfig:      app.cmdGenConfig,
		service.CmdRegister:       app.cmdRegister,
		service.CmdProfile:        app.cmdProfile,
		service.CmdChangeKeys:     app.cmdChangeKeys,
		service.CmdUpdatePassword: app.cmdUpdatePassword,
		service.CmdGet:            app.cmdGet,
		service.CmdSave:           app.cmdSave,
		service.CmdDelete:         app.cmdDelete,
		service.CmdList:           app.cmdList,
	}, app.logger)

	return nil
}

func (app *App) getCmd(cmd *config.AppCommands) (service.Command, error) {
	if cmd.GenConfig {
		return service.CmdGenConfig, nil
	}
	if cmd.Register {
		return service.CmdRegister, nil
	}

	if cmd.Profile {
		return service.CmdProfile, nil
	}
	if cmd.GenKeys {
		return service.CmdChangeKeys, nil
	}
	if cmd.ChangePassword {
		return service.CmdUpdatePassword, nil
	}

	if cmd.Get {
		return service.CmdGet, nil
	}
	if cmd.Put {
		return service.CmdSave, nil
	}
	if cmd.Remove {
		return service.CmdDelete, nil
	}
	if cmd.List {
		return service.CmdList, nil
	}

	return "", errs.NewAppCommonError("no command applied", nil)
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
