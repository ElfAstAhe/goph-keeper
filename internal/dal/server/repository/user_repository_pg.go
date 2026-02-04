package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/model"
	irepo "github.com/ElfAstAhe/goph-keeper/internal/bll/server/repository"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

// все sql запросы в рамках репозитория
const (
	sqlGet      string = "select id, username, password_hash, private_key, public_key, active, person, e_mail from users where id = $1"
	sqlGetByKey string = "select id, username, password_hash, private_key, public_key, active, person, e_mail from users where username = $1"
	sqlCreate   string = "insert into users(id, username, password_hash, private_key, public_key, active, person, e_mail) values ($1, $2, $3, $4, $5, $6, $7, $8)"
	sqlChange   string = `update users set
                 password_hash = $2,
                 private_key = $3,
                 public_key = $4,
                 active = $5,
                 person = $6,
                 e_mail = $7
where id = $1`
	sqlRemove string = "delete from users where id = $1"
)

type UserRepositoryPg struct {
	db               utils.DB
	dataCipherHelper *utils.CipherHelper
	userDataRepo     irepo.UserDataRepository
}

func NewUserRepositoryPg(db utils.DB, dataCipherHelper *utils.CipherHelper, userDataRepo irepo.UserDataRepository) *UserRepositoryPg {
	return &UserRepositoryPg{
		db:               db,
		dataCipherHelper: dataCipherHelper,
		userDataRepo:     userDataRepo,
	}
}

func (urp *UserRepositoryPg) Get(ctx context.Context, id string) (*model.User, error) {
	res, err := urp.internalGetSingle(ctx, sqlGet, id)
	if err != nil {
		return nil, apperrs.NewDalRepositoryError("user repo", "get by id", err)
	}

	return res, nil
}

func (urp *UserRepositoryPg) GetByKey(ctx context.Context, key *model.UserKey) (*model.User, error) {
	res, err := urp.internalGetSingle(ctx, sqlGetByKey, key.Username)
	if err != nil {
		return nil, apperrs.NewDalRepositoryError("user repo", "get by key", err)
	}

	return res, nil
}

func (urp *UserRepositoryPg) internalGetSingle(ctx context.Context, sqlReq string, params ...any) (*model.User, error) {
	row := urp.db.GetDB().QueryRowContext(ctx, sqlReq, params...)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	entity := model.NewEmptyUser()
	err := row.Scan(&entity.ID, &entity.Key.Username, &entity.PasswordHash, &entity.PrivateKey, &entity.Active, &entity.Person, &entity.EMail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return entity, nil
}

func (urp *UserRepositoryPg) internalGetMulti(ctx context.Context, sql string, params ...any) ([]model.User, error) {
	return nil, errs.NewAppCommonError("not implemented", nil)
}

func (urp *UserRepositoryPg) Create(ctx context.Context, user *model.User) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (urp *UserRepositoryPg) Change(ctx context.Context, user *model.User) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (urp *UserRepositoryPg) Remove(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (urp *UserRepositoryPg) ListUserData(ctx context.Context, id string) ([]*model.UserData, error) {
	return urp.userDataRepo.ListByOwner(ctx, id)
}
