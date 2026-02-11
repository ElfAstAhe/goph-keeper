package service

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

func (ch *CmdHandlerImpl) cmdUpdatePassword(ctx context.Context, options *config.AppOptions) error {
	// валидация
	err := ch.validateUpdatePassword(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("UpdatePassword validation failed.", err)
	}
	// предвариловка
	err = ch.beforeCmd(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("UpdatePassword beforeCmd", err)
	}
	// выполнение
	// аутентификация
	err = ch.client.Login(ctx, ch.settings.GetConfig().Username, ch.settings.GetConfig().EncryptedPassword)
	if err != nil {
		return errs.NewAppCommonError("Error logging in", err)
	}

	ch.log.Info("Successfully logged in")

	err = ch.client.UpdatePassword(ctx, options.OldPassword, options.NewPassword)
	if err != nil {
		return errs.NewAppCommonError("Error updating password", err)
	}

	err = ch.postUpdatePassword(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("Error post update password", err)
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

func (ch *CmdHandlerImpl) validateUpdatePassword(ctx context.Context, options *config.AppOptions) error {
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
	if options.OldPassword == "" {
		return errs.NewAppInvalidArgumentError("oldPassword", "empty")
	}
	if options.NewPassword == "" {
		return errs.NewAppInvalidArgumentError("newPassword", "empty")
	}

	return nil
}

func (ch *CmdHandlerImpl) postUpdatePassword(ctx context.Context, options *config.AppOptions) error {
	var err error
	// шифруем пароль
	ch.settings.GetConfig().EncryptedPassword, err = ch.encryptPassword(options.NewPassword, ch.settings.GetConfig().PublicKey)
	if err != nil {
		return errs.NewAppConfigError("encrypt password", err)
	}

	return nil
}
