package service

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

func (ch *CmdHandlerImpl) cmdProfile(ctx context.Context, options *config.AppOptions) error {
	// валидация
	err := ch.validateProfile(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("Profile validation failed.", err)
	}
	// предвариловка
	err = ch.beforeCmd(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("Profile beforeCmd", err)
	}
	// выполнение
	// аутентификация
	err = ch.client.Login(ctx, ch.settings.GetConfig().Username, ch.settings.GetConfig().EncryptedPassword)
	if err != nil {
		return errs.NewAppCommonError("Error logging in", err)
	}

	ch.log.Info("Successfully logged in")

	// получаем профиль
	res, err := ch.client.GetProfile(ctx)
	if err != nil {
		return errs.NewAppCommonError("Error getting profile", err)
	}

	ch.log.Info("Successfully got profile")

	// выводим профиль
	err = ch.logProfile(ctx, res)
	if err != nil {
		return errs.NewAppCommonError("Error report profile", err)
	}

	// постобработка
	err = ch.postProfile(ctx, options, res)

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

func (ch *CmdHandlerImpl) validateProfile(ctx context.Context, options *config.AppOptions) error {
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

	return nil
}

func (ch *CmdHandlerImpl) logProfile(ctx context.Context, result *dto.UserDto) error {
	ch.log.Warnf("[profile]\n username [%s]\n person [%s]\n e-mail [%s]\n public key [%s]", result.Username, result.Person, result.EMail, result.PublicKey)

	return nil
}

func (ch *CmdHandlerImpl) postProfile(ctx context.Context, options *config.AppOptions, res *dto.UserDto) error {
	ch.settings.GetConfig().PublicKey = res.PublicKey

	return nil
}
