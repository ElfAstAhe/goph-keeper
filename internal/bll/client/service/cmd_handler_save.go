package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	"github.com/ElfAstAhe/goph-keeper/pkg/client/rest"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

func (ch *CmdHandlerImpl) cmdSave(ctx context.Context, options *config.AppOptions) error {
	// валидация
	err := ch.validateSave(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("Save validation failed.", err)
	}
	// предвариловка
	err = ch.beforeCmd(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("Save beforeCmd", err)
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
		} else if ok && clientErr.StatusCode == http.StatusGone {
			ch.log.Warnf("user data already removed, please, change name, data kind [%s], name [%s]", options.DataKind, options.Name)

			return nil
		} else {
			return errs.NewAppCommonError("Error getting data", err)
		}
	}

	// отправляем данные
	res, err = ch.saveUserData(ctx, res, options)
	if err != nil {
		return errs.NewAppCommonError("Error saving data", err)
	}

	ch.log.Info("Successfully saved data")

	return nil
}

func (ch *CmdHandlerImpl) validateSave(ctx context.Context, options *config.AppOptions) error {
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
	if options.DataKind == dto.UserDataKindBinary {
		if options.Path == "" {
			return errs.NewAppInvalidArgumentError("path", "empty")
		}
		if _, err := os.Stat(options.Path); os.IsNotExist(err) {
			return errs.NewAppConfigItemError("file not found", err)
		}
	}

	return nil
}

func (ch *CmdHandlerImpl) saveUserData(ctx context.Context, res *dto.UserDataDto, options *config.AppOptions) (*dto.UserDataDto, error) {
	if options.DataKind == dto.UserDataKindBinary {
		return ch.saveBinaryData(ctx, res, options)
	}

	return ch.saveTextData(ctx, res, options)
}

func (ch *CmdHandlerImpl) saveBinaryData(ctx context.Context, res *dto.UserDataDto, options *config.AppOptions) (*dto.UserDataDto, error) {
	var err error
	f, err := os.OpenFile(options.Path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, errs.NewAppCommonError("Error opening file", err)
	}
	defer f.Close()

	res.DataKind = options.DataKind
	res.Name = options.Name
	res.TextData = filepath.Base(options.Path)
	res.BinaryData, err = io.ReadAll(f)
	res.CreatedAt = time.Now()
	res.ModifiedAt = time.Now()

	res, err = ch.client.Save(ctx, res)
	if err != nil {
		return nil, errs.NewAppCommonError("Error saving data", err)
	}

	return res, nil
}

func (ch *CmdHandlerImpl) saveTextData(ctx context.Context, res *dto.UserDataDto, options *config.AppOptions) (*dto.UserDataDto, error) {
	res.DataKind = options.DataKind
	res.Name = options.Name
	res.TextData = options.Data
	if options.Path != "" {
		content, err := os.ReadFile(options.Path)
		if err != nil {
			return nil, errs.NewAppCommonError("Error reading file", err)
		}
		res.TextData = string(content)
	}
	res.CreatedAt = time.Now()
	res.ModifiedAt = time.Now()

	res, err := ch.client.Save(ctx, res)
	if err != nil {
		return nil, errs.NewAppCommonError("Error saving data", err)
	}

	return res, nil
}
