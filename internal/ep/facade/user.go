package facade

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
)

type UserFacade interface {
	GetProfile(ctx context.Context) (*dto.UserDto, error)
	UpdatePassword(ctx context.Context, changePassword *dto.UpdatePasswordDto) error
	ChangeKeys(ctx context.Context) (*dto.ChangeKeysResultDto, error)
}
