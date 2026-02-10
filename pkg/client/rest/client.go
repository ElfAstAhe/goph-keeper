package rest

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
)

type GophKeeperClient interface {
	Login(ctx context.Context, username, encryptedPassword string) error
	Register(ctx context.Context, username, password, person, eMail string) (*dto.RegisterResultDto, error)
	GetProfile(ctx context.Context) (*dto.UserDto, error)
	ChangeKeys(ctx context.Context) (*dto.ChangeKeysResultDto, error)
	UpdatePassword(ctx context.Context, oldPassword, newPassword string) error
	Get(ctx context.Context, id string) (*dto.UserDataDto, error)
	GetByKey(ctx context.Context, dataKind string, name string) (*dto.UserDataDto, error)
	ListAll(ctx context.Context) ([]*dto.UserDataDto, error)
	Save(ctx context.Context, data *dto.UserDataDto) (*dto.UserDataDto, error)
	Delete(ctx context.Context, id string) error
}
