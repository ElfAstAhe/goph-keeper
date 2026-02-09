package service

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/model"
)

type UserDataService interface {
	Get(ctx context.Context, userID string, key *model.UserDataKey) (*model.UserData, error)
	Save(ctx context.Context, userID string, userData *model.UserData) (*model.UserData, error)
	ListAll(ctx context.Context, userID string) ([]*model.UserDataShort, error)
	Delete(ctx context.Context, userID string, ID string) error
}
