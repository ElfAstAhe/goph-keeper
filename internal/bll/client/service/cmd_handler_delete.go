package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	"github.com/ElfAstAhe/goph-keeper/pkg/client/rest"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

func (ch *CmdHandlerImpl) cmdDelete(ctx context.Context, options *config.AppOptions) error {
	// валидация
	err := ch.validateDelete(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("Delete validation failed.", err)
	}
	// предвариловка
	err = ch.beforeCmd(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("Delete beforeCmd", err)
	}
	// выполнение
	// аутентификация
	err = ch.client.Login(ctx, ch.settings.GetConfig().Username, ch.settings.GetConfig().EncryptedPassword)
	if err != nil {
		return errs.NewAppCommonError("Error logging in", err)
	}

	ch.log.Info("Successfully logged in")

	// подгружаем данные
	res, err := ch.client.GetByKey(ctx, options.DataKind, options.Name)
	if err != nil {
		var clientErr *rest.ClientError
		ok := errors.As(err, &clientErr)
		if ok && clientErr.StatusCode == http.StatusNotFound {
			ch.log.Warnf("user data not found, data kind [%s], name [%s]", options.DataKind, options.Name)

			res = dto.NewEmptyUserDataDto()
		} else {
			return errs.NewAppCommonError("Error getting data", err)
		}
	}

	err = ch.client.Delete(ctx, res.ID)
	if err != nil {
		return errs.NewAppCommonError("Error deleting data", err)
	}

	ch.log.Info("Successfully deleted data")

	return nil
}

func (ch *CmdHandlerImpl) validateDelete(ctx context.Context, options *config.AppOptions) error {
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
