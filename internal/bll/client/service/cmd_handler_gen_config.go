package service

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

func (ch *CmdHandlerImpl) cmdGenConfig(ctx context.Context, options *config.AppOptions) error {
	err := ch.settings.SaveFile()
	if err != nil {
		return errs.NewAppCommonError("save config file", err)
	}

	ch.log.Infof("config file saved at [%s]", options.ConfigPath)

	return nil
}
