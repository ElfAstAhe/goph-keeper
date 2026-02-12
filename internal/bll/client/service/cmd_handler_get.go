package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	"github.com/ElfAstAhe/goph-keeper/pkg/client/rest"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

func (ch *CmdHandlerImpl) cmdGet(ctx context.Context, options *config.AppOptions) error {
	// валидация
	err := ch.validateGet(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("Get validation failed.", err)
	}
	// предвариловка
	err = ch.beforeCmd(ctx, options)
	if err != nil {
		return errs.NewAppCommonError("Get beforeCmd", err)
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
		var clientErr *rest.ClientError
		ok := errors.As(err, &clientErr)
		if ok && clientErr.StatusCode == http.StatusNotFound {
			ch.log.Warnf("user data not found, data kind [%s], name [%s]", options.DataKind, options.Name)

			return nil
		}
		if ok && clientErr.StatusCode == http.StatusGone {
			ch.log.Warnf("user data removed, data kind [%s], name [%s]", options.DataKind, options.Name)

			return nil
		}

		return errs.NewAppCommonError("Error getting data", err)
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

	ch.log.Info("Successfully post get response")

	return nil
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
	var err error
	// пост обработка только для типа binary, сохранение файла
	if opts.DataKind == dto.UserDataKindBinary {
		if opts.Path == "" {
			opts.Path, err = ch.buildDefaultBinaryDataPath(result)
			if err != nil {
				return errs.NewAppCommonError("error build binary data path", err)
			}
		}

		f, err := os.OpenFile(opts.Path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			return errs.NewAppCommonError("error open file", err)
		}
		defer f.Close()

		buf := bytes.NewBuffer(result.BinaryData)
		_, err = io.Copy(f, buf)
		if err != nil {
			return errs.NewAppCommonError(fmt.Sprintf("error write file at [%s]", opts.Path), err)
		}

		ch.log.Infof("binary data stored at [%s]", opts.Path)
	}

	return nil
}

func (ch *CmdHandlerImpl) logGet(ctx context.Context, res *dto.UserDataDto) error {
	ch.log.Infof("[user data]\n data kind: %s\n name: %s\n created at: [%v]\n modified at [%v]\n text data: %s", res.DataKind, res.Name, res.CreatedAt, res.ModifiedAt, res.TextData)
	if res.DataKind == dto.UserDataKindBinary {
		ch.log.Infof("binary data length [%d]", len(res.BinaryData))
	}

	return nil
}

func (ch *CmdHandlerImpl) buildDefaultBinaryDataPath(result *dto.UserDataDto) (string, error) {
	// 1. Получаем путь к запущенному бинарнику (например, /home/user/goph-keeper/bin/server)
	exePath, err := os.Executable()
	if err != nil {
		return "", errs.NewAppConfigError("get executable path", err)
	}

	// 2. Берем директорию бинарника (/home/user/goph-keeper/bin)
	exeDir := filepath.Dir(exePath)

	// 4. Собираем путь (например, /home/user/goph-keeper/config/config.yaml)
	configPath := filepath.Join(exeDir, fmt.Sprintf("%s-%s-%s-binary.gk", result.ID, result.DataKind, result.Name))

	return configPath, nil
}
