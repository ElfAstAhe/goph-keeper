package bootstrap

import (
	"github.com/ElfAstAhe/goph-keeper/internal/app/db"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/ElfAstAhe/goph-keeper/pkg/logger"
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
	var cipherKey []byte = make([]byte, 0, 32)
	// key cipher
	app.keyCipher = utils.NewSHA256Cipher()
	// prepare correct cipher key
	cipherKey, err = app.keyCipher.Encrypt([]byte(app.conf.CipherKey))
	if err != nil {
		return errs.NewAppCommonError("init helpers build correct cipher key error", err)
	}
	// cipher
	app.cipher, err = utils.NewAesGcmCipher(cipherKey)
	if err != nil {
		return errs.NewAppCommonError("init helpers cipher error", err)
	}
	// cipher helper
	app.cipherHelper = utils.NewCipherHelper(app.cipher)
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
	// ToDo: implement

	return nil
}

func (app *App) initDependencies() error {
	// ToDo: implement

	return nil
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
