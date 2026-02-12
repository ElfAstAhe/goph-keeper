package service

import (
	"context"
	"fmt"

	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

type CmdService interface {
	Execute(ctx context.Context, cmd *config.AppCommands, opts *config.AppOptions) error
}

type CmdServiceImpl struct {
	cmdHandler CmdHandler
}

func NewCmdServiceImpl(cmdHandler CmdHandler) *CmdServiceImpl {
	return &CmdServiceImpl{
		cmdHandler: cmdHandler,
	}
}

func (cs *CmdServiceImpl) Execute(ctx context.Context, cmd *config.AppCommands, opts *config.AppOptions) error {
	// валидация
	if err := cs.validate(cmd); err != nil {
		return err
	}
	// преобразование
	command, err := cs.incomeToCommand(cmd)
	if err != nil {
		return errs.NewAppCommonError("income to command transform", err)
	}
	// выполнение
	err = cs.cmdHandler.Process(ctx, command, opts)
	if err != nil {
		return errs.NewAppCommonError("command processing", err)
	}

	return nil
}

func (cs *CmdServiceImpl) validate(cmd *config.AppCommands) error {
	if cmd == nil {
		return errs.NewAppCommonError("empty income command", nil)
	}
	if err := cmd.Validate(); err != nil {
		return errs.NewAppCommonError("validate income command", err)
	}

	return nil
}

func (cs *CmdServiceImpl) incomeToCommand(income *config.AppCommands) (Command, error) {
	if income == nil {
		return "", errs.NewAppCommonError("empty income command", nil)
	}

	if income.GenConfig {
		return CmdGenConfig, nil
	}
	if income.Register {
		return CmdRegister, nil
	}

	if income.Profile {
		return CmdProfile, nil
	}
	if income.GenKeys {
		return CmdChangeKeys, nil
	}
	if income.ChangePassword {
		return CmdUpdatePassword, nil
	}

	if income.Get {
		return CmdGet, nil
	}
	if income.Put {
		return CmdSave, nil
	}
	if income.Remove {
		return CmdDelete, nil
	}
	if income.List {
		return CmdList, nil
	}

	return "", errs.NewAppCommonError(fmt.Sprintf("unknown command: %v", income), nil)
}
