package mapper

import (
	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/model"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

func RegisterDtoToUser(source *dto.RegisterDto) (*model.User, error) {
	if source == nil {
		return nil, apperrs.NewEpMappingError("RegisterDtoToUser", "RegisterDto", "User", "empty source", nil)
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
		return nil, apperrs.NewEpMappingError("UserToRegisterResultDto", "User", "RegisterResultDto", "nil source", nil)
	}

	res := dto.NewRegisterResultDto(src.PublicKey)

	return res, nil
}

func UserToUserDto(src *model.User) (*dto.UserDto, error) {
	if src == nil {
		return nil, apperrs.NewEpMappingError("UserToUserDto", "User", "UserDto", "nil source", nil)
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

func UserDataToUserDataDto(src *model.UserData) (*dto.UserDataDto, error) {
	if src == nil {
		return nil, apperrs.NewEpMappingError("UserDataToUserDataDto", "UserData", "UserDataDto", "nil source", nil)
	}

	res := dto.NewEmptyUserDataDto()

	res.ID = src.ID
	res.DataKind = src.Key.DataKind
	res.Name = src.Key.Name
	res.TextData = src.TextData
	res.BinaryData = make([]byte, len(src.BinaryData))
	copy(res.BinaryData, src.BinaryData)
	res.CreatedAt = src.CreatedAt
	res.ModifiedAt = src.ModifiedAt

	return res, nil
}

func UserDataDtoToUserData(src *dto.UserDataDto) (*model.UserData, error) {
	if src == nil {
		return nil, apperrs.NewEpMappingError("UserDataDtoToUserData", "UserDataDto", "UserData", "nil source", nil)
	}

	res := model.NewEmptyUserData()

	res.ID = src.ID
	res.Key = model.NewUserDataKey(src.Name, src.DataKind)
	res.TextData = src.TextData
	res.BinaryData = make([]byte, len(src.BinaryData))
	copy(res.BinaryData, src.BinaryData)
	res.CreatedAt = src.CreatedAt
	res.ModifiedAt = src.ModifiedAt

	return res, nil
}

func UserDataShortToUserDataDto(src *model.UserDataShort) (*dto.UserDataDto, error) {
	if src == nil {
		return nil, apperrs.NewEpMappingError("UserDataShortToUserDataDto", "UserDataShort", "UserData", "nil source", nil)
	}

	res := dto.NewEmptyUserDataDto()

	res.ID = src.ID
	res.DataKind = src.Key.DataKind
	res.Name = src.Key.Name
	res.CreatedAt = src.CreatedAt
	res.ModifiedAt = src.ModifiedAt

	return res, nil
}

func UserDataShortListToUserDataDto(src []*model.UserDataShort) ([]*dto.UserDataDto, error) {
	if len(src) == 0 {
		return make([]*dto.UserDataDto, 0), nil
	}

	res := make([]*dto.UserDataDto, 0, len(src))
	for _, udMdl := range src {
		udDto, err := UserDataShortToUserDataDto(udMdl)
		if err != nil {
			return nil, err
		}

		res = append(res, udDto)
	}

	return res, nil
}
