package facade

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/model"
	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/service"
	"github.com/ElfAstAhe/goph-keeper/internal/ep/mapper"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

type UserDataFacadeImpl struct {
	userDataService service.UserDataService
	authHelper      *utils.AuthHelper
}

func NewUserDataFacadeImpl(userDataService service.UserDataService, authHelper *utils.AuthHelper) *UserDataFacadeImpl {
	return &UserDataFacadeImpl{
		userDataService: userDataService,
		authHelper:      authHelper,
	}
}

func (udf *UserDataFacadeImpl) Get(ctx context.Context, ID string) (*dto.UserDataDto, error) {
	userInfo, err := udf.authHelper.UserInfoFromContext(ctx)
	if err != nil {
		return nil, errs.NewAuthForbiddenError("user info from context", err)
	}

	userData, err := udf.userDataService.Get(ctx, userInfo.UserID(), ID)
	if err != nil {
		return nil, err
	}

	return mapper.UserDataToUserDataDto(userData)
}

func (udf *UserDataFacadeImpl) GetByKey(ctx context.Context, dataKind, name string) (*dto.UserDataDto, error) {
	userInfo, err := udf.authHelper.UserInfoFromContext(ctx)
	if err != nil {
		return nil, errs.NewAuthForbiddenError("user info from context", err)
	}

	userData, err := udf.userDataService.GetByKey(ctx, userInfo.UserID(), model.NewUserDataKey(name, dataKind))
	if err != nil {
		return nil, err
	}

	return mapper.UserDataToUserDataDto(userData)
}

func (udf *UserDataFacadeImpl) ListAll(ctx context.Context) ([]*dto.UserDataDto, error) {
	userInfo, err := udf.authHelper.UserInfoFromContext(ctx)
	if err != nil {
		return nil, errs.NewAuthForbiddenError("user info from context", err)
	}

	userData, err := udf.userDataService.ListAll(ctx, userInfo.UserID())
	if err != nil {
		return nil, err
	}

	return mapper.UserDataShortListToUserDataDto(userData)
}

func (udf *UserDataFacadeImpl) Create(ctx context.Context, data *dto.UserDataDto) (*dto.UserDataDto, error) {
	userInfo, err := udf.authHelper.UserInfoFromContext(ctx)
	if err != nil {
		return nil, errs.NewAuthForbiddenError("user info from context", err)
	}

	userData, err := mapper.UserDataDtoToUserData(data)
	if err != nil {
		return nil, err
	}

	userData.ID = ""

	userData, err = udf.userDataService.Save(ctx, userInfo.UserID(), userData)
	if err != nil {
		return nil, err
	}

	return mapper.UserDataToUserDataDto(userData)
}

func (udf *UserDataFacadeImpl) Change(ctx context.Context, data *dto.UserDataDto) (*dto.UserDataDto, error) {
	userInfo, err := udf.authHelper.UserInfoFromContext(ctx)
	if err != nil {
		return nil, errs.NewAuthForbiddenError("user info from context", err)
	}

	userData, err := mapper.UserDataDtoToUserData(data)
	if err != nil {
		return nil, err
	}

	userData, err = udf.userDataService.Save(ctx, userInfo.UserID(), userData)
	if err != nil {
		return nil, err
	}

	return mapper.UserDataToUserDataDto(userData)
}

func (udf *UserDataFacadeImpl) Remove(ctx context.Context, ID string) error {
	userInfo, err := udf.authHelper.UserInfoFromContext(ctx)
	if err != nil {
		return errs.NewAuthForbiddenError("user info from context", err)
	}

	err = udf.userDataService.Delete(ctx, userInfo.UserID(), ID)
	if err != nil {
		return err
	}

	return nil
}
