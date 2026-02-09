package facade

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/ep/dto"
)

type UserDataFacade interface {
	Get(ctx context.Context, ID string) (*dto.UserDataDto, error)
	GetByKey(ctx context.Context, dataKind, name string) (*dto.UserDataDto, error)

	Create(ctx context.Context, data *dto.UserDataDto) (*dto.UserDataDto, error)
	Change(ctx context.Context, data *dto.UserDataDto) (*dto.UserDataDto, error)

	Remove(ctx context.Context, ID string) error
}
