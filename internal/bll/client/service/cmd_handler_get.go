package service

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

func (ch *CmdHandlerImpl) cmdGet(ctx context.Context, options *config.AppOptions) error {
	// валидация
	err := ch.validateGet(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("CHangeKeys validation failed.", err)
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

	// получаем данные
	res, err := ch.client.GetByKey(ctx, options.DataKind, options.Name)
	if err != nil {
		return errs.NewAppCommonError("GetByKey Error", err)
	}

	ch.log.Info("Successfully get by key")

	err = ch.logGet(ctx, res)
	if err != nil {
		return errs.NewAppCommonError("error report user data", err)
	}

	err = ch.postGet(ctx, options, res)
	if err != nil {
		return errs.NewAppCommonError("PostGet Error", err)
	}

	return errs.NewAppCommonError("not implemented", nil)
}

func (ch *CmdHandlerImpl) validateGet(ctx context.Context, options *config.AppOptions) error {
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
	if options.DataKind == "" {
		return errs.NewAppInvalidArgumentError("data_kind", "empty")
	}
	if options.Name == "" {
		return errs.NewAppInvalidArgumentError("name", "empty")
	}

	return nil
}

func (ch *CmdHandlerImpl) postGet(ctx context.Context, opts *config.AppOptions, result *dto.UserDataDto) error {
	// пост обработка только для типа binary, сохранение файла
	if opts.DataKind != dto.UserDataKindBinary {
		return nil
	}

	return nil
}

func (ch *CmdHandlerImpl) logGet(ctx context.Context, res *dto.UserDataDto) error {
	ch.log.Infof("[user data]\n data kind [%s]\n name [%s]\n created at [%v] modified at [%v]", res.DataKind, res.Name, res.CreatedAt, res.ModifiedAt)

}
