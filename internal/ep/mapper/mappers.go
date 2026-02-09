package mapper

import (
	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/model"
	"github.com/ElfAstAhe/goph-keeper/internal/ep/dto"
	"github.com/ElfAstAhe/goph-keeper/internal/err"
)

func RegisterDtoToUser(source *dto.RegisterDto) (*model.User, error) {
	if source == nil {
		return nil, err.NewEpMappingError("RegisterDtoToUser", "RegisterDto", "User", "empty source", nil)
	}

	res := model.NewEmptyUser()

	res.Key = model.NewUserKey(source.Username)
	res.Active = true
	res.Deleted = false
	res.Person = source.Person
	res.EMail = source.EMail

	return res, nil
}

func UserToRegisterResultDto(src *model.User) (*dto.RegisterResultDto, error) {
	if src == nil {
		return nil, err.NewEpMappingError("UserToRegisterResultDto", "User", "RegisterResultDto", "nil source", nil)
	}

	res := dto.NewRegisterResultDto(src.PublicKey)

	return res, nil
}

func UserToUserDto(src *model.User) (*dto.UserDto, error) {
	if src == nil {
		return nil, err.NewEpMappingError("UserToUserDto", "User", "UserDto", "nil source", nil)
	}

	res := dto.NewEmptyUserDto()

	res.ID = src.ID
	res.Username = src.Key.Username
	res.PasswordHash = src.PasswordHash
	res.PublicKey = src.PublicKey
	res.Active = src.Active
	res.Person = src.Person
	res.EMail = src.EMail

	return res, nil
}
