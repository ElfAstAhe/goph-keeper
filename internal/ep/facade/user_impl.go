package facade

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/service"
	"github.com/ElfAstAhe/goph-keeper/internal/ep/mapper"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

type UserFacadeImpl struct {
	userService service.UserService
	authHelper  *utils.AuthHelper
}

func NewUserFacadeImpl(userService service.UserService, authHelper *utils.AuthHelper) *UserFacadeImpl {
	return &UserFacadeImpl{
		userService: userService,
		authHelper:  authHelper,
	}
}

func (uf *UserFacadeImpl) GetProfile(ctx context.Context) (*dto.UserDto, error) {
	userInfo, err := uf.authHelper.UserInfoFromContext(ctx)
	if err != nil {
		return nil, errs.NewAuthForbiddenError("extract user info from context", err)
	}
	user, err := uf.userService.GetProfile(ctx, userInfo.UserID())
	if err != nil {
		return nil, err
	}

	return mapper.UserToUserDto(user)
}

func (uf *UserFacadeImpl) UpdatePassword(ctx context.Context, changePassword *dto.UpdatePasswordDto) error {
	userInfo, err := uf.authHelper.UserInfoFromContext(ctx)
	if err != nil {
		return errs.NewAuthForbiddenError("extract user info from context", err)
	}

	err = uf.userService.UpdatePassword(ctx, userInfo.UserID(), changePassword.NewPassword, changePassword.OldPassword)
	if err != nil {
		return err
	}

	return nil
}

func (uf *UserFacadeImpl) ChangeKeys(ctx context.Context) (*dto.ChangeKeysResultDto, error) {
	userInfo, err := uf.authHelper.UserInfoFromContext(ctx)
	if err != nil {
		return nil, errs.NewAuthForbiddenError("extract user info from context", err)
	}

	publicKey, err := uf.userService.ChangeKeys(ctx, userInfo.UserID())
	if err != nil {
		return nil, err
	}

	return dto.NewChangeKeysResultDto(publicKey), nil
}
