package repository

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/model"
)

type UserRepository interface {
	Get(ctx context.Context, id string) (*model.User, error)
	GetByKey(ctx context.Context, key *model.UserKey) (*model.User, error)

	Create(ctx context.Context, user *model.User) (*model.User, error)
	Change(ctx context.Context, user *model.User) (*model.User, error)

	Remove(ctx context.Context, id string) error

	ListUserData(ctx context.Context, id string) ([]*model.UserDataShort, error)
}
