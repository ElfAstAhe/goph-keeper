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
		return nil, err
	}

	return urp.afterGet(res)
}

func (urp *UserRepositoryPg) GetByKey(ctx context.Context, key *model.UserKey) (*model.User, error) {
	res, err := urp.internalGetSingle(ctx, sqlUserGetByKey, key.Username)
	if err != nil {
		return nil, err
	}

	return urp.afterGet(res)
}

func (urp *UserRepositoryPg) afterGet(user *model.User) (*model.User, error) {
	var err error = nil
	// расшифровываем данные
	user.PrivateKey = urp.dataCipherHelper.DecryptString(user.PrivateKey)
	user.PublicKey = urp.dataCipherHelper.DecryptString(user.PublicKey)

	if user.Deleted {
		err = apperrs.NewDalSoftDeletedError("User", fmt.Sprintf("ID: [%s], Key: [%s]", user.ID, user.Key))
	}

	return user, err
}

func (urp *UserRepositoryPg) internalGetSingle(ctx context.Context, sqlReq string, params ...any) (*model.User, error) {
	row := urp.db.GetDB().QueryRowContext(ctx, sqlReq, params...)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, apperrs.NewDalCommonError("UserRepo.internalGetSingle", "fetch row", row.Err())
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
			return nil, apperrs.NewDalNotFoundError("UserRepo.internalGetSingle", "user not found", err)
		}

		return nil, apperrs.NewDalCommonError("UserRepo.internalGetSingle", "scan row", err)
	}

	return entity, nil
}

func (urp *UserRepositoryPg) internalGetMulti(ctx context.Context, sqlReq string, params ...any) ([]model.User, error) {
	return nil, errs.NewAppCommonError("not implemented", nil)
}

func (urp *UserRepositoryPg) Create(ctx context.Context, user *model.User) (res *model.User, err error) {
	// валидируем
	if err = urp.validateCreate(user); err != nil {
		return nil, apperrs.NewDalValidateError("User", "validate create", err)
	}

	// подготавливаем
	if err = urp.beforeCreate(user); err != nil {
		return nil, apperrs.NewDalCommonError("UserRepo.Create", "before create", err)
	}

	// выполнение
	err = urp.db.GetHelper().RunInTx(ctx, urp.db.GetDB(), func(tx *sql.Tx) error {
		// стейтмент
		stmt, err := tx.PrepareContext(ctx, sqlUserCreate)
		if err != nil {
			return apperrs.NewDalCommonError("UserRepo.Create", "prepare stmt", err)
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx,
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
			if urp.db.GetHelper().IsUniqueViolation(err) {
				return apperrs.NewDalAlreadyExistsError("User", user.Key.Username, err)
			}

			return apperrs.NewDalCommonError("UserRepo.Create", "exec stmt", err)
		}

		return nil
	})
	if err != nil {
		return nil, apperrs.NewDalCommonError("UserRepo.Create", "run in transaction", err)
	}

	return urp.afterGet(user)
}

func (urp *UserRepositoryPg) validateCreate(user *model.User) error {
	if user == nil {
		return errs.NewAppInvalidArgumentError("user", "nil user")
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

func (urp *UserRepositoryPg) Change(ctx context.Context, user *model.User) (res *model.User, err error) {
	// валидируем
	if err = urp.validateChange(user); err != nil {
		return nil, apperrs.NewDalValidateError("User", "validate change", err)
	}
	// подготавливаем
	if err = urp.beforeChange(user); err != nil {
		return nil, apperrs.NewDalCommonError("UserRepo.Change", "before change", err)
	}

	// сохраняем
	err = urp.db.GetHelper().RunInTx(ctx, urp.db.GetDB(), func(tx *sql.Tx) error {
		// стейтмент
		stmt, err := tx.PrepareContext(ctx, sqlUserChange)
		if err != nil {
			return apperrs.NewDalCommonError("UserRepo.Change", "prepare stmt", err)
		}
		defer stmt.Close()

		err = urp.db.GetHelper().ExecStmt(ctx, stmt, func(repErr error) (string, string, error) {
			return "User", user.GetID(), repErr
		},
			user.ID,
			user.PasswordHash,
			user.PrivateKey,
			user.PublicKey,
			user.Active,
			user.Deleted,
			user.Person,
			user.EMail,
		)
		if err != nil {
			return apperrs.NewDalCommonError("UserRepo.Change", "exec stmt", err)
		}

		return nil
	})
	if err != nil {
		return nil, apperrs.NewDalCommonError("UserRepo.Change", "run in transaction", err)
	}

	return urp.afterGet(user)
}

func (urp *UserRepositoryPg) validateChange(user *model.User) error {
	if user == nil {
		return errs.NewAppInvalidArgumentError("user", "nil user")
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

func (urp *UserRepositoryPg) Remove(ctx context.Context, id string) (err error) {
	// сохраняем
	err = urp.db.GetHelper().RunInTx(ctx, urp.db.GetDB(), func(tx *sql.Tx) error {
		// стейтмент
		stmt, err := tx.PrepareContext(ctx, sqlUserRemove)
		if err != nil {
			return apperrs.NewDalCommonError("UserRepo.Remove", "prepare stmt", err)
		}
		defer stmt.Close()

		err = urp.db.GetHelper().ExecStmt(ctx, stmt, func(repErr error) (string, string, error) {
			return "User", id, repErr
		}, id)
		if err != nil {
			return apperrs.NewDalCommonError("UserRepo.Remove", "exec stmt", err)
		}

		return nil
	})
	if err != nil {
		return apperrs.NewDalCommonError("UserRepo.Remove", "run in transaction", err)
	}

	return nil
}

func (urp *UserRepositoryPg) ListUserData(ctx context.Context, id string) ([]*model.UserDataShort, error) {
	return urp.userDataRepo.ListAllByOwner(ctx, id)
}
