package service

import (
	"context"
	"fmt"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

func (ch *CmdHandlerImpl) cmdList(ctx context.Context, options *config.AppOptions) error {
	// валидация
	err := ch.validateList(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("List validation failed.", err)
	}
	// предвариловка
	err = ch.beforeCmd(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("List beforeCmd", err)
	}
	// выполнение
	// аутентификация
	err = ch.client.Login(ctx, ch.settings.GetConfig().Username, ch.settings.GetConfig().EncryptedPassword)
	if err != nil {
		return errs.NewAppCommonError("Error logging in", err)
	}

	ch.log.Info("Successfully logged in")

	res, err := ch.client.ListAll(ctx)
	if err != nil {
		return errs.NewAppCommonError("Error listing all", err)
	}

	err = ch.logList(ctx, res)
	if err != nil {
		return errs.NewAppCommonError("Error report all user data", err)
	}

	ch.log.Info("Successfully list all user data")

	return nil
}

func (ch *CmdHandlerImpl) validateList(ctx context.Context, options *config.AppOptions) error {
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

func (ch *CmdHandlerImpl) logList(ctx context.Context, res []*dto.UserDataDto) error {
	ch.log.Info(ch.buildTableHeader())
	for _, item := range res {
		ch.log.Info(ch.buildTableString(item))
	}

	return nil
}

func (ch *CmdHandlerImpl) buildTableHeader() string {
	return fmt.Sprintf("%-36s | %-15s | %-16s | %-16s | %s", "ID", "DATA KIND", "Cr AT", "Md AT", "NAME")
}

func (ch *CmdHandlerImpl) buildTableString(data *dto.UserDataDto) string {
	return fmt.Sprintf("%-36s | %-15s | %-16s | %-16s | %s", data.ID, data.DataKind, data.CreatedAt.Format("02.01.2006 15:04"), data.ModifiedAt.Format("02.01.2006 15:04"), data.Name)
}
