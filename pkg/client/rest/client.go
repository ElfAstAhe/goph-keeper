package rest

import (
	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
)

type GophKeeperClient interface {
	Login(username, encryptedPassword string) error
	Register(username, password, person, eMail string) (*dto.RegisterResultDto, error)
	GetProfile() (*dto.UserDto, error)
	ChangeKeys() (*dto.ChangeKeysResultDto, error)
	UpdatePassword(oldPassword, newPassword string) error
	Get(id string) (*dto.UserDataDto, error)
	GetByKey(dataKind string, name string) (*dto.UserDataDto, error)
	ListAll() ([]*dto.UserDataDto, error)
	Save(data *dto.UserDataDto) (*dto.UserDataDto, error)
	Delete(id string) error
}
