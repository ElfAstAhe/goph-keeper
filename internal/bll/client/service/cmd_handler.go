package service

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	"github.com/ElfAstAhe/goph-keeper/pkg/client/rest"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

type CmdHandlerFunc func(ctx context.Context, options *config.AppOptions) error

type CmdHandler interface {
	Process(ctx context.Context, command Command, opts *config.AppOptions) error
}

type CmdHandlerImpl struct {
	conf     *config.AppConfig
	client   rest.GophKeeperClient
	handlers map[Command]CmdHandlerFunc
}

func NewCmdHandlerImpl(client rest.GophKeeperClient, conf *config.AppConfig) *CmdHandlerImpl {
	res := &CmdHandlerImpl{
		client: client,
		conf:   conf,
	}

	return res.init()
}

func (cp *CmdHandlerImpl) Process(ctx context.Context, command Command, opts *config.AppOptions) error {

	return errs.NewAppCommonError("not implemented", nil)
}

func (cp *CmdHandlerImpl) init() *CmdHandlerImpl {
	cp.handlers = make(map[Command]CmdHandlerFunc)

	cp.handlers[CmdGenConfig] = cp.cmdGenConfig
	cp.handlers[CmdRegister] = cp.cmdRegister

	cp.handlers[CmdProfile] = cp.cmdProfile
	cp.handlers[CmdChangeKeys] = cp.cmdChangeKeys
	cp.handlers[CmdUpdatePassword] = cp.cmdUpdatePassword

	cp.handlers[CmdGet] = cp.cmdGet
	cp.handlers[CmdSave] = cp.cmdSave
	cp.handlers[CmdDelete] = cp.cmdDelete
	cp.handlers[CmdList] = cp.cmdList

	return cp
}

func (cp *CmdHandlerImpl) cmdGenConfig(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}

func (cp *CmdHandlerImpl) cmdRegister(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}

func (cp *CmdHandlerImpl) cmdProfile(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}

func (cp *CmdHandlerImpl) cmdChangeKeys(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}

func (cp *CmdHandlerImpl) cmdUpdatePassword(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}

func (cp *CmdHandlerImpl) cmdGet(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}

func (cp *CmdHandlerImpl) cmdSave(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}

func (cp *CmdHandlerImpl) cmdDelete(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}

func (cp *CmdHandlerImpl) cmdList(ctx context.Context, options *config.AppOptions) error {
	// ToDo: implement

	return errs.NewAppCommonError("not implemented", nil)
}
