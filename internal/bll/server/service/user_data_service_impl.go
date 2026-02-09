package service

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/model"
	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/repository"
)

type UserDataServiceImpl struct {
	userDataRepo repository.UserDataRepository
}

func NewUserDataService(userDataRepo repository.UserDataRepository) *UserDataServiceImpl {
	return &UserDataServiceImpl{
		userDataRepo: userDataRepo,
	}
}

func (uds *UserDataServiceImpl) Get(ctx context.Context, userID string, ID string) (*model.UserData, error) {
	res, err := uds.userDataRepo.Get(ctx, userID, ID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (uds *UserDataServiceImpl) GetByKey(ctx context.Context, userID string, key *model.UserDataKey) (*model.UserData, error) {
	res, err := uds.userDataRepo.GetByKey(ctx, userID, key)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (uds *UserDataServiceImpl) Save(ctx context.Context, userID string, userData *model.UserData) (*model.UserData, error) {
	var err error
	var res *model.UserData
	if userData.ID == "" {
		res, err = uds.userDataRepo.Create(ctx, userID, userData)
	} else {
		res, err = uds.userDataRepo.Change(ctx, userID, userData)
	}
	if err != nil {
		return nil, err
	}

	// сбрасываем тяжёлые данные
	res.TextData = ""
	res.BinaryData = make([]byte, 0)

	return res, nil
}

func (uds *UserDataServiceImpl) ListAll(ctx context.Context, userID string) ([]*model.UserDataShort, error) {
	return uds.userDataRepo.ListAllByOwner(ctx, userID)
}

func (uds *UserDataServiceImpl) Delete(ctx context.Context, userID string, ID string) error {
	return uds.userDataRepo.Remove(ctx, userID, ID)
}
