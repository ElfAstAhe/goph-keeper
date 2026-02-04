package bootstrap

import (
	"github.com/ElfAstAhe/goph-keeper/internal/app/db"
	"github.com/ElfAstAhe/goph-keeper/internal/dal/server/repository"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/ElfAstAhe/goph-keeper/pkg/logger"
	migrations "github.com/ElfAstAhe/goph-keeper/pkg/migrations/goose"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

func (app *App) loadConfig() error {
	log := app.log.GetLogger("bootstrap load config")
	if err := app.conf.Load(); err != nil {
		return errs.NewAppCommonError("load config error", err)
	}

	log.Infof("config FINAL: [%+v]", app.conf)

	if err := app.conf.Validate(); err != nil {
		return errs.NewAppConfigError("validate config error", err)
	}

	return nil
}

func (app *App) initLogger() error {
	appLogger, err := logger.NewZapLogger("INFO", "")
	if err != nil {
		return errs.NewAppCommonError("init logger error", err)
	}

	app.log = appLogger

	return nil
}

func (app *App) initHelpers() error {
	var err error
	var cipherKey = make([]byte, 0, 32)
	// key cipher
	app.keyCipher = utils.NewSHA256Cipher()
	// prepare correct cipher key
	cipherKey, err = app.keyCipher.Encrypt([]byte(app.conf.CipherKey))
	if err != nil {
		return errs.NewAppCommonError("init helpers build correct cipher key error", err)
	}
	// data cipher
	app.dataCipher, err = utils.NewAesGcmCipher(cipherKey)
	if err != nil {
		return errs.NewAppCommonError("init helpers cipher error", err)
	}
	// cipher helper
	app.dataCipherHelper = utils.NewCipherHelper(app.dataCipher)
	// jwt helper
	app.jwtHelper = utils.NewDefaultJWTHelper(app.conf.JWTConfig.SecretKey)
	// auth helper
	app.authHelper = utils.NewDefaultAuthHelperEx(app.jwtHelper)

	return nil
}

func (app *App) initDatabase() error {
	var err error
	app.db, err = db.NewPostgresDB(app.conf.DatabaseConfig)
	if err != nil {
		return errs.NewAppCommonError("init database error", err)
	}

	return nil
}

func (app *App) migrateDatabase() error {
	migrator, err := migrations.NewGooseDBMigrator(app.ctx, app.db.GetDB(), app.log)
	if err != nil {
		return errs.NewAppCommonError("create migrator", err)
	}
	if err := migrator.Initialize(); err != nil {
		return errs.NewAppCommonError("init migrator", err)
	}
	if err := migrator.Up(); err != nil {
		return errs.NewAppCommonError("migrator up", err)
	}

	return nil
}

func (app *App) initDependencies() error {
	var err error

	// repositories
	app.userDataRepo = repository.NewUserDataRepositoryPg(app.db, app.dataCipherHelper)
	app.userRepo = repository.NewUserRepositoryPg(app.db, app.dataCipherHelper, app.userDataRepo)

	// services

	return err
}

func (app *App) initStartupServices() error {
	// ToDo: implement

	return nil
}

func (app *App) initHTTPRouter() error {
	// ToDo: implement

	return nil
}

func (app *App) initHTTPServer() error {
	// ToDo: implement

	return nil
}

func (app *App) initGRPCServer() error {
	// ToDo: implement

	return nil
}
