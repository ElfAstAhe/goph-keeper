package facade

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/ep/dto"
)

type AuthFacade interface {
	Login(ctx context.Context, login *dto.LoginDto) (*dto.LoginResultDto, error)
	Register(ctx context.Context, register *dto.RegisterDto) (*dto.RegisterResultDto, error)
}
