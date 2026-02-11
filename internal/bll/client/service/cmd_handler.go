package service

import (
	"context"
	"fmt"

	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	"github.com/ElfAstAhe/goph-keeper/pkg/client/rest"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/ElfAstAhe/goph-keeper/pkg/logger"
)

type CmdHandlerFunc func(ctx context.Context, options *config.AppOptions) error

type CmdHandler interface {
	Process(ctx context.Context, command Command, opts *config.AppOptions) error
}

type CmdHandlerImpl struct {
	settings *config.AppSettings
	client   rest.GophKeeperClient
	handlers map[Command]CmdHandlerFunc
	log      logger.Logger
}

func NewCmdHandlerImpl(client rest.GophKeeperClient, settings *config.AppSettings, logger logger.Logger) *CmdHandlerImpl {
	res := &CmdHandlerImpl{
		client:   client,
		settings: settings,
		log:      logger.GetLogger("cmd handler"),
	}

	return res.init()
}

func (ch *CmdHandlerImpl) Process(ctx context.Context, command Command, opts *config.AppOptions) error {
	handler, err := ch.getHandlerFund(command)
	if err != nil {
		return errs.NewAppCommonError("get handler", err)
	}

	err = handler(ctx, opts)
	if err != nil {
		return errs.NewAppCommonError("handler process", err)
	}

	return nil
}

func (ch *CmdHandlerImpl) init() *CmdHandlerImpl {
	ch.handlers = make(map[Command]CmdHandlerFunc)

	ch.handlers[CmdGenConfig] = ch.cmdGenConfig
	ch.handlers[CmdRegister] = ch.cmdRegister

	ch.handlers[CmdProfile] = ch.cmdProfile
	ch.handlers[CmdChangeKeys] = ch.cmdChangeKeys
	ch.handlers[CmdUpdatePassword] = ch.cmdUpdatePassword

	ch.handlers[CmdGet] = ch.cmdGet
	ch.handlers[CmdSave] = ch.cmdSave
	ch.handlers[CmdDelete] = ch.cmdDelete
	ch.handlers[CmdList] = ch.cmdList

	return ch
}

func (ch *CmdHandlerImpl) getHandlerFund(cmd Command) (CmdHandlerFunc, error) {
	res, ok := ch.handlers[cmd]
	if !ok {
		return nil, errs.NewAppCommonError(fmt.Sprintf("not found handler for command [%s]", cmd), nil)
	}

	return res, nil
}

func (ch *CmdHandlerImpl) cmdGenConfig(ctx context.Context, options *config.AppOptions) error {
	err := ch.settings.SaveFile()
	if err != nil {
		return errs.NewAppCommonError("save config file", err)
	}

	ch.log.Infof("config file saved at [%s]", options.ConfigPath)

	return nil
}

func (ch *CmdHandlerImpl) cmdRegister(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}

func (ch *CmdHandlerImpl) cmdProfile(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}

func (ch *CmdHandlerImpl) cmdChangeKeys(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}

func (ch *CmdHandlerImpl) cmdUpdatePassword(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}

func (ch *CmdHandlerImpl) cmdGet(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}

func (ch *CmdHandlerImpl) cmdSave(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}

func (ch *CmdHandlerImpl) cmdDelete(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}

func (ch *CmdHandlerImpl) cmdList(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}
