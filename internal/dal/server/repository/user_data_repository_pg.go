package repository

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/model"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

type UserDataRepositoryPg struct {
	db               utils.DB
	dataCipherHelper *utils.CipherHelper
}

func NewUserDataRepositoryPg(db utils.DB, dataCipherHelper *utils.CipherHelper) *UserDataRepositoryPg {
	return &UserDataRepositoryPg{
		db:               db,
		dataCipherHelper: dataCipherHelper,
	}
}

func (u UserDataRepositoryPg) Get(ctx context.Context, id string) (*model.UserData, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserDataRepositoryPg) GetByKey(ctx context.Context, key *model.UserDataKey) (*model.UserData, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserDataRepositoryPg) Create(ctx context.Context, user *model.UserData) (*model.UserData, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserDataRepositoryPg) Change(ctx context.Context, user *model.UserData) (*model.UserData, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserDataRepositoryPg) Remove(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (u UserDataRepositoryPg) GetByOwner(ctx context.Context, userID string) (*model.UserData, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserDataRepositoryPg) ListByOwner(ctx context.Context, userID string) ([]*model.UserData, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserDataRepositoryPg) RemoveAllByOwner(ctx context.Context, userID string) error {
	//TODO implement me
	panic("implement me")
}
