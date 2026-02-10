package service

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	"github.com/ElfAstAhe/goph-keeper/pkg/logger"
)

type ProcessorFunc func(ctx context.Context, conf *config.AppConfig, options *config.AppOptions) error

type CommandProcessor struct {
	cmdProcessors map[Command]ProcessorFunc
	log           logger.Logger
}

func NewCommandProcessor(cmdProcessors map[Command]ProcessorFunc, logger logger.Logger) *CommandProcessor {
	return &CommandProcessor{
		cmdProcessors: cmdProcessors,
		log:           logger.GetLogger("CommandProcessor"),
	}
}

func (cp *CommandProcessor) Process(command Command, conf *config.AppConfig, options *config.AppOptions) error {
	// ToDo: implement

	return nil
}
