package repository

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/model"
)

type UserDataRepository interface {
	Get(ctx context.Context, id string) (*model.UserData, error)
	GetByKey(ctx context.Context, userID string, key *model.UserDataKey) (*model.UserData, error)

	Create(ctx context.Context, userID string, user *model.UserData) (*model.UserData, error)
	Change(ctx context.Context, userID string, user *model.UserData) (*model.UserData, error)

	Remove(ctx context.Context, id string) error

	GetByOwner(ctx context.Context, userID string) (*model.UserData, error)
	ListByOwner(ctx context.Context, userID string) ([]*model.UserData, error)
	RemoveAllByOwner(ctx context.Context, userID string) error
}
