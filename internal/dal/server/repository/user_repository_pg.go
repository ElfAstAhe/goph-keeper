package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/model"
	irepo "github.com/ElfAstAhe/goph-keeper/internal/bll/server/repository"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

// все sql запросы в рамках репозитория
const (
	sqlUserGet string = `
select
    id,
    username,
    password_hash,
    private_key,
    public_key,
    active,
    deleted,
    person,
    e_mail
from
    users
where
    id = $1`

	sqlUserGetByKey string = `
select
    id,
    username,
    password_hash,
    private_key,
    public_key,
    active,
    deleted,
    person,
    e_mail
from
    users
where
    username = $1`

	sqlUserCreate string = `
insert into users(
    id,
    username,
    password_hash,
    private_key,
    public_key,
    active,
    deleted,
    person,
    e_mail
)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	sqlUserChange string = `
update
    users
set
    password_hash = $2,
    private_key = $3,
    public_key = $4,
    active = $5,
    deleted = $6,
    person = $7,
    e_mail = $8
where
    id = $1`

	sqlUserRemove string = `
update
    users
set
    deleted = true
where
    id = $1`
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
	res, err := urp.internalGetSingle(ctx, sqlUserGet, id)
	if err != nil {
		return nil, apperrs.NewDalRepositoryError("UserRepo.Get", "get by id", err)
	}

	return urp.afterGet(res)
}

func (urp *UserRepositoryPg) GetByKey(ctx context.Context, key *model.UserKey) (*model.User, error) {
	res, err := urp.internalGetSingle(ctx, sqlUserGetByKey, key.Username)
	if err != nil {
		return nil, apperrs.NewDalRepositoryError("UserRepo.GetByKey", "get by key", err)
	}

	return urp.afterGet(res)
}

func (urp *UserRepositoryPg) afterGet(user *model.User) (*model.User, error) {
	var err error = nil
	// расшифровываем данные
	user.PrivateKey = urp.dataCipherHelper.DecryptString(user.PrivateKey)
	user.PublicKey = urp.dataCipherHelper.DecryptString(user.PublicKey)

	if user.Deleted {
		err = apperrs.NewBllModelSoftDeletedError("user")
	}

	return user, err
}

func (urp *UserRepositoryPg) internalGetSingle(ctx context.Context, sqlReq string, params ...any) (*model.User, error) {
	row := urp.db.GetDB().QueryRowContext(ctx, sqlReq, params...)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	entity := model.NewEmptyUser()
	err := row.Scan(
		&entity.ID,
		&entity.Key.Username,
		&entity.PasswordHash,
		&entity.PrivateKey,
		&entity.PublicKey,
		&entity.Active,
		&entity.Deleted,
		&entity.Person,
		&entity.EMail,
	)
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

func (urp *UserRepositoryPg) Create(ctx context.Context, user *model.User) (res *model.User, err error) {
	// валидируем
	if err = urp.validateCreate(user); err != nil {
		return nil, apperrs.NewDalRepositoryError("UserRepo.Create", "validate create", err)
	}
	// подготавливаем
	if err := user.BeforeCreate(); err != nil {
		return nil, apperrs.NewDalRepositoryError("UserRepo.Create", "before create", err)
	}

	// сохраняем
	// транзакция
	tx, err := urp.db.GetDB().Begin()
	if err != nil {
		return nil, apperrs.NewDalRepositoryError("UserRepo.Create", "begin transaction", err)
	}
	defer func() {
		if r := recover(); r != nil {
			err = tx.Rollback()

			recoveryErr, ok := r.(error)
			if !ok {
				recoveryErr = fmt.Errorf("%v", r)
			}
			err = apperrs.NewDalRepositoryError("UserRepo.Create", "rollback transaction", recoveryErr)

			return
		}

		err = tx.Commit()
	}()

	// стейтмент
	stmt, err := tx.PrepareContext(ctx, sqlUserCreate)
	if err != nil {
		return nil, apperrs.NewDalRepositoryError("UserRepo.Create", "create sql statement", err)
	}
	defer stmt.Close()

	res, err = urp.createStmt(ctx, stmt, user)
	if err != nil {
		return nil, apperrs.NewDalRepositoryError("UserRepo.Create", "insert data", err)
	}

	return res, err
}

func (urp *UserRepositoryPg) createStmt(ctx context.Context, stmt *sql.Stmt, user *model.User) (*model.User, error) {
	_, err := stmt.ExecContext(ctx,
		user.ID,
		user.Key.Username,
		user.PasswordHash,
		user.PrivateKey,
		user.PublicKey,
		user.Active,
		user.Deleted,
		user.Person,
		user.EMail,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (urp *UserRepositoryPg) validateCreate(user *model.User) error {
	if user == nil {
		return apperrs.NewDalRepositoryError("user repo", "validate create user", nil)
	}

	return user.ValidateCreate()
}

func (urp *UserRepositoryPg) beforeCreate(user *model.User) error {
	if err := user.BeforeCreate(); err != nil {
		return err
	}

	// шифруем данные (проверка на уже зашифрованные данные внутри хелпера)
	user.PrivateKey = urp.dataCipherHelper.EncryptString(user.PrivateKey)
	user.PublicKey = urp.dataCipherHelper.EncryptString(user.PublicKey)

	return nil
}

func (urp *UserRepositoryPg) Change(ctx context.Context, user *model.User) (*model.User, error) {
	// ToDo: implement

	return nil, nil
}

func (urp *UserRepositoryPg) validateChange(user *model.User) error {
	if user == nil {
		return apperrs.NewDalRepositoryError("user repo", "validate change user", nil)
	}

	return user.ValidateChange()
}

func (urp *UserRepositoryPg) beforeChange(user *model.User) error {
	if err := user.BeforeChange(); err != nil {
		return err
	}

	// шифруем данные
	user.PrivateKey = urp.dataCipherHelper.EncryptString(user.PrivateKey)
	user.PublicKey = urp.dataCipherHelper.EncryptString(user.PublicKey)

	return nil
}

func (urp *UserRepositoryPg) Remove(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (urp *UserRepositoryPg) ListUserData(ctx context.Context, id string) ([]*model.UserData, error) {
	return urp.userDataRepo.ListByOwner(ctx, id)
}
