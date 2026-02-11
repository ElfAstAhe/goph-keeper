package service

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

func (ch *CmdHandlerImpl) cmdChangeKeys(ctx context.Context, options *config.AppOptions) error {
	// валидация
	err := ch.validateChangeKeys(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("ChangeKeys validation failed.", err)
	}
	// предвариловка
	err = ch.beforeCmd(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("ChangeKeys beforeCmd", err)
	}
	// выполнение
	// аутентификация
	err = ch.client.Login(ctx, ch.settings.GetConfig().Username, ch.settings.GetConfig().EncryptedPassword)
	if err != nil {
		return errs.NewAppCommonError("Error logging in", err)
	}

	ch.log.Info("Successfully logged in")

	res, err := ch.client.ChangeKeys(ctx)
	if err != nil {
		return errs.NewAppCommonError("Error change keys", err)
	}

	ch.log.Info("Successfully changed keys")

	err = ch.postChangeKeys(ctx, options, res)
	if err != nil {
		return errs.NewAppCommonError("error post change keys", err)
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

func (ch *CmdHandlerImpl) validateChangeKeys(ctx context.Context, options *config.AppOptions) error {
	if ch.settings.GetConfig().Address == "" {
		return errs.NewAppInvalidArgumentError("address", "empty")
	}
	if ch.settings.GetConfig().Username == "" {
		return errs.NewAppInvalidArgumentError("username", "empty")
	}
	if ch.settings.GetConfig().EncryptedPassword == "" &&
		options.Password == "" {
		return errs.NewAppInvalidArgumentError("password", "empty")
	}
	if options.Password == "" {
		return errs.NewAppInvalidArgumentError("password", "empty")
	}

	return nil
}

func (ch *CmdHandlerImpl) postChangeKeys(ctx context.Context, options *config.AppOptions, result *dto.ChangeKeysResultDto) error {
	var err error
	// шифруем пароль
	ch.settings.GetConfig().EncryptedPassword, err = ch.encryptPassword(options.Password, result.PublicKey)
	if err != nil {
		return errs.NewAppConfigError("encrypt password", err)
	}
	ch.settings.GetConfig().PublicKey = result.PublicKey

	return nil
}
