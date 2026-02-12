package service

import (
	"context"
	"fmt"

	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	"github.com/ElfAstAhe/goph-keeper/pkg/client/rest"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/ElfAstAhe/goph-keeper/pkg/logger"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

type CmdHandlerFunc func(ctx context.Context, options *config.AppOptions) error

type CmdHandler interface {
	Process(ctx context.Context, command Command, opts *config.AppOptions) error
}

type CmdHandlerImpl struct {
	keysHelper *utils.RSAKeysHelper
	settings   *config.AppSettings
	client     rest.GophKeeperClient
	handlers   map[Command]CmdHandlerFunc
	log        logger.Logger
}

func NewCmdHandlerImpl(keysHelper *utils.RSAKeysHelper, client rest.GophKeeperClient, settings *config.AppSettings, logger logger.Logger) *CmdHandlerImpl {
	res := &CmdHandlerImpl{
		keysHelper: keysHelper,
		client:     client,
		settings:   settings,
		log:        logger.GetLogger("cmd handler"),
	}

	return res.init()
}

func (ch *CmdHandlerImpl) Process(ctx context.Context, command Command, opts *config.AppOptions) error {
	handler, err := ch.getHandlerFunc(command)
	if err != nil {
		return errs.NewAppCommonError("get handler", err)
	}

	ch.log.Infof("command [%s] processing", command)

	err = handler(ctx, opts)
	if err != nil {
		return errs.NewAppCommonError("handler process", err)
	}

	ch.log.Infof("command [%s] processed", command)

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

func (ch *CmdHandlerImpl) getHandlerFunc(cmd Command) (CmdHandlerFunc, error) {
	res, ok := ch.handlers[cmd]
	if !ok {
		return nil, errs.NewAppCommonError(fmt.Sprintf("not found handler for command [%s]", cmd), nil)
	}

	return res, nil
}

func (ch *CmdHandlerImpl) encryptPassword(password string, publicKey string) (string, error) {
	pubKey, err := ch.keysHelper.ParsePublicKey(publicKey)
	if err != nil {
		return "", errs.NewAppConfigError("parse public key", err)
	}

	encryptedPassword, err := ch.keysHelper.EncryptString(password, pubKey)
	if err != nil {
		return "", errs.NewAppConfigError("encrypt password", err)
	}

	return encryptedPassword, nil
}

func (ch *CmdHandlerImpl) beforeCmd(ctx context.Context, opts *config.AppOptions) error {
	if opts == nil {
		return errs.NewAppCommonError("options cannot be nil", nil)
	}
	if opts.Password != "" {
		encryptedPassword, err := ch.encryptPassword(opts.Password, ch.settings.GetConfig().PublicKey)
		if err != nil {
			return errs.NewAppCommonError("encrypt password", err)
		}
		ch.settings.GetConfig().EncryptedPassword = encryptedPassword
	}

	return nil
}
