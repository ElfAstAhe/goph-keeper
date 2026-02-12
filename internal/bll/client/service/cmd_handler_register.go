package service

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

func (ch *CmdHandlerImpl) cmdRegister(ctx context.Context, options *config.AppOptions) error {
	// валидация
	err := ch.validateRegister(options)
	if err != nil {
		return errs.NewAppCommonError("validate register", err)
	}
	// выполнение
	res, err := ch.client.Register(ctx, options.Username, options.Password, options.Person, options.EMail)
	if err != nil {
		return errs.NewAppCommonError("register", err)
	}

	ch.log.Infof("Registration completed successfully")

	// постобработка
	err = ch.postRegister(ctx, options, res)
	if err != nil {
		return errs.NewAppCommonError("post register", err)
	}

	// сохраняем конфиг
	if !options.NotStoreConfig {
		err = ch.settings.SaveFile()
		if err != nil {
			return errs.NewAppCommonError("save config file", err)
		}

		ch.log.Infof("config file saved at [%s]", options.ConfigPath)
	}

	return nil
}

func (ch *CmdHandlerImpl) validateRegister(opts *config.AppOptions) error {
	if ch.settings.GetConfig().Address == "" {
		return errs.NewAppInvalidArgumentError("address", "empty")
	}
	if opts == nil {
		return errs.NewAppInvalidArgumentError("opts", "empty")
	}
	if opts.Username == "" {
		return errs.NewAppInvalidArgumentError("opts.Username", "empty")
	}
	if opts.Password == "" {
		return errs.NewAppInvalidArgumentError("opts.Password", "empty")
	}

	return nil
}

func (ch *CmdHandlerImpl) postRegister(ctx context.Context, opts *config.AppOptions, result *dto.RegisterResultDto) error {
	var err error
	// шифруем пароль
	ch.settings.GetConfig().EncryptedPassword, err = ch.encryptPassword(opts.Password, result.PublicKey)
	if err != nil {
		return errs.NewAppConfigError("encrypt password", err)
	}
	ch.settings.GetConfig().PublicKey = result.PublicKey

	return nil
}
