package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/model"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

const (
	sqlUserDataGet string = `
select
    id,
    name,
    kind,
    text_data,
    binary_data,
    created_at,
    modified_at,
    deleted
from
    user_data
where
    id = $1
`

	sqlUserDataGetByKey string = `
select
    id,
    name,
    kind,
    text_data,
    binary_data,
    created_at,
    modified_at,
    deleted
from
    user_data
where
    user_id = $1
and name = $2
and kind = $3
`

	sqlUserDataCreate string = `
insert into user_data (
    id,
    user_id,
    name,
    kind,
    text_data,
    binary_data,
    created_at,
    modified_at,
    deleted
)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`

	sqlUserDataChange string = `
update
    user_data
set
    name = $3,
    kind = $4,
    text_data = $5,
    binary_data = $6,
    created_at = $7,
    modified_at = $8,
    deleted = $9
where
    user_id = $2
and id = $1
`

	sqlUserDataRemove string = `
update
    user_data
set
    deleted = true,
    modified_at = now()
where
    id = $1
`

	sqlUserDataRemoveAllByOwner string = `
update
    user_data
set
    deleted = true,
    modified_at = now()
where
    deleted = false
and user_id = $1
`
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

func (udrp *UserDataRepositoryPg) Get(ctx context.Context, id string) (*model.UserData, error) {
	res, err := udrp.internalGetSingle(ctx, sqlUserDataGet, id)
	if err != nil {
		return nil, apperrs.NewDalRepositoryError("UserDataRepo.Get", "get by id", err)
	}

	return udrp.afterGet(res)
}

func (udrp *UserDataRepositoryPg) GetByKey(ctx context.Context, userID string, key *model.UserDataKey) (*model.UserData, error) {
	res, err := udrp.internalGetSingle(ctx, sqlUserDataGetByKey, userID, key.Name, key.DataKind)
	if err != nil {
		return nil, apperrs.NewDalRepositoryError("UserDataRepo.GetByKey", "get by key", err)
	}

	return udrp.afterGet(res)
}

func (udrp *UserDataRepositoryPg) internalGetSingle(ctx context.Context, sqlReq string, params ...any) (*model.UserData, error) {
	row := udrp.db.GetDB().QueryRowContext(ctx, sqlReq, params...)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	entity := model.NewEmptyUserData()
	err := row.Scan(
		&entity.ID,
		&entity.Key.Name,
		&entity.Key.DataKind,
		&entity.TextData,
		&entity.BinaryData,
		&entity.CreatedAt,
		&entity.ModifiedAt,
		&entity.Deleted,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return entity, nil
}

func (udrp *UserDataRepositoryPg) afterGet(res *model.UserData) (*model.UserData, error) {
	var err error = nil
	// расшифровываем данные
	res.TextData = udrp.dataCipherHelper.DecryptString(res.TextData)
	res.BinaryData = udrp.dataCipherHelper.DecryptBinary(res.BinaryData)

	if res.Deleted {
		err = apperrs.NewBllModelSoftDeletedError("UserData")
	}

	return res, err
}

func (udrp *UserDataRepositoryPg) Create(ctx context.Context, userID string, userData *model.UserData) (res *model.UserData, err error) {
	// валидируем
	if err = udrp.validateCreate(userID, userData); err != nil {
		return nil, apperrs.NewDalRepositoryError("UserDataRepo.Create", "validate create", err)
	}
	// подготавливаем
	if err = udrp.beforeCreate(userData); err != nil {
		return nil, apperrs.NewDalRepositoryError("UserDataRepo.Create", "before create", err)
	}

	// сохраняем
	// транзакция
	tx, err := udrp.db.GetDB().Begin()
	if err != nil {
		return nil, apperrs.NewDalRepositoryError("UserDataRepo.Create", "begin transaction", err)
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback() // Откатываем в любом случае

			// Превращаем панику в читаемую ошибку для логов
			var recoveryErr error
			if e, ok := r.(error); ok {
				recoveryErr = e
			} else {
				recoveryErr = fmt.Errorf("%v", r)
			}

			err = apperrs.NewDalRepositoryError("UserDataRepo.Create", "panic recovery", recoveryErr)
		} else if err != nil {
			_ = tx.Rollback() // Откат при ошибке бизнеса/БД
		} else {
			err = tx.Commit() // Фиксация
			if err != nil {
				err = apperrs.NewDalRepositoryError("UserDataRepo.Create", "commit", err)
			}
		}
	}()

	// стейтмент
	stmt, err := tx.PrepareContext(ctx, sqlUserDataCreate)
	if err != nil {
		return nil, apperrs.NewDalRepositoryError("UserDataRepo.Create", "create sql statement", err)
	}
	defer stmt.Close()

	err = udrp.execStmt(ctx, stmt,
		userData.ID,
		userID,
		userData.Key.Name,
		userData.Key.DataKind,
		userData.TextData,
		userData.BinaryData,
		userData.CreatedAt,
		userData.ModifiedAt,
		userData.Deleted,
	)
	if err != nil {
		return nil, apperrs.NewDalRepositoryError("UserDataRepo.Create", "insert data", err)
	}

	return udrp.afterGet(userData)
}

func (udrp *UserDataRepositoryPg) execStmt(ctx context.Context, stmt *sql.Stmt, params ...any) error {
	_, err := stmt.ExecContext(ctx, params...)

	if err != nil {
		return err
	}

	return nil
}

func (udrp *UserDataRepositoryPg) validateCreate(userID string, userData *model.UserData) error {
	if userData == nil {
		return apperrs.NewDalRepositoryError("UserDataRepo.validateCreate", "nil user data instance", nil)
	}
	if userID == "" {
		return apperrs.NewDalRepositoryError("UserDataRepo.validateCreate", "empty user ID", nil)
	}

	return userData.ValidateCreate()
}

func (udrp *UserDataRepositoryPg) beforeCreate(userData *model.UserData) error {
	if err := userData.BeforeCreate(); err != nil {
		return err
	}

	// шифруем
	userData.TextData = udrp.dataCipherHelper.EncryptString(userData.TextData)
	userData.BinaryData = udrp.dataCipherHelper.EncryptBinary(userData.BinaryData)

	return nil
}

func (udrp *UserDataRepositoryPg) Change(ctx context.Context, userID string, userData *model.UserData) (res *model.UserData, err error) {
	// валидируем
	if err = udrp.validateChange(userID, userData); err != nil {
		return nil, apperrs.NewDalRepositoryError("UserDataRepo.Change", "validate change", err)
	}
	// подготавливаем
	if err = udrp.beforeChange(userData); err != nil {
		return nil, apperrs.NewDalRepositoryError("UserDataRepo.Change", "before change", err)
	}

	// сохраняем
	// транзакция
	tx, err := udrp.db.GetDB().Begin()
	if err != nil {
		return nil, apperrs.NewDalRepositoryError("UserDataRepo.Change", "begin transaction", err)
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback() // Откатываем в любом случае

			// Превращаем панику в читаемую ошибку для логов
			var recoveryErr error
			if e, ok := r.(error); ok {
				recoveryErr = e
			} else {
				recoveryErr = fmt.Errorf("%v", r)
			}

			err = apperrs.NewDalRepositoryError("UserDataRepo.Change", "panic recovery", recoveryErr)
		} else if err != nil {
			_ = tx.Rollback() // Откат при ошибке бизнеса/БД
		} else {
			err = tx.Commit() // Фиксация
			if err != nil {
				err = apperrs.NewDalRepositoryError("UserDataRepo.Change", "commit", err)
			}
		}
	}()

	// стейтмент
	stmt, err := tx.PrepareContext(ctx, sqlUserDataChange)
	if err != nil {
		return nil, apperrs.NewDalRepositoryError("UserDataRepo.Change", "change sql statement", err)
	}
	defer stmt.Close()

	err = udrp.execStmt(ctx, stmt,
		userData.ID,
		userID,
		userData.Key.Name,
		userData.Key.DataKind,
		userData.TextData,
		userData.BinaryData,
		userData.CreatedAt,
		userData.ModifiedAt,
		userData.Deleted,
	)
	if err != nil {
		return nil, apperrs.NewDalRepositoryError("UserDataRepo.Change", "update data", err)
	}

	return udrp.afterGet(userData)
}

func (udrp *UserDataRepositoryPg) validateChange(userID string, userData *model.UserData) error {
	if userData == nil {
		return apperrs.NewDalRepositoryError("UserDataRepo.validateChange", "nil user data instance", nil)
	}
	if userID == "" {
		return apperrs.NewDalRepositoryError("UserDataRepo.validateChange", "empty user ID", nil)
	}

	return userData.ValidateChange()
}

func (udrp *UserDataRepositoryPg) beforeChange(userData *model.UserData) error {
	if err := userData.BeforeChange(); err != nil {
		return err
	}

	// шифруем
	userData.TextData = udrp.dataCipherHelper.EncryptString(userData.TextData)
	userData.BinaryData = udrp.dataCipherHelper.EncryptBinary(userData.BinaryData)

	return nil
}

func (udrp *UserDataRepositoryPg) Remove(ctx context.Context, id string) error {
	// сохраняем
	// транзакция
	tx, err := udrp.db.GetDB().Begin()
	if err != nil {
		return apperrs.NewDalRepositoryError("UserDataRepo.Remove", "begin transaction", err)
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback() // Откатываем в любом случае

			// Превращаем панику в читаемую ошибку для логов
			var recoveryErr error
			if e, ok := r.(error); ok {
				recoveryErr = e
			} else {
				recoveryErr = fmt.Errorf("%v", r)
			}

			err = apperrs.NewDalRepositoryError("UserDataRepo.Remove", "panic recovery", recoveryErr)
		} else if err != nil {
			_ = tx.Rollback() // Откат при ошибке бизнеса/БД
		} else {
			err = tx.Commit() // Фиксация
			if err != nil {
				err = apperrs.NewDalRepositoryError("UserDataRepo.Remove", "commit", err)
			}
		}
	}()
	// стейтмент
	stmt, err := tx.PrepareContext(ctx, sqlUserDataRemove)
	if err != nil {
		return apperrs.NewDalRepositoryError("UserDataRepo.Remove", "remove sql statement", err)
	}
	defer stmt.Close()

	err = udrp.execStmt(ctx, stmt, id)
	if err != nil {
		return apperrs.NewDalRepositoryError("UserDataRepo.Remove", "update data", err)
	}

	return nil
}

func (udrp *UserDataRepositoryPg) ListAllByOwner(ctx context.Context, userID string) ([]*model.UserData, error) {
	//TODO implement me
	panic("implement me")
}

func (udrp *UserDataRepositoryPg) internalGetMulti(ctx context.Context, sqlReq string, params ...any) ([]model.UserData, error) {
	return nil, errs.NewAppCommonError("not implemented", nil)
}

func (udrp *UserDataRepositoryPg) RemoveAllByOwner(ctx context.Context, userID string) (err error) {
	// сохраняем
	// транзакция
	tx, err := udrp.db.GetDB().Begin()
	if err != nil {
		return apperrs.NewDalRepositoryError("UserDataRepo.RemoveAllByOwner", "begin transaction", err)
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback() // Откатываем в любом случае

			// Превращаем панику в читаемую ошибку для логов
			var recoveryErr error
			if e, ok := r.(error); ok {
				recoveryErr = e
			} else {
				recoveryErr = fmt.Errorf("%v", r)
			}

			err = apperrs.NewDalRepositoryError("UserDataRepo.RemoveAllByOwner", "panic recovery", recoveryErr)
		} else if err != nil {
			_ = tx.Rollback() // Откат при ошибке бизнеса/БД
		} else {
			err = tx.Commit() // Фиксация
			if err != nil {
				err = apperrs.NewDalRepositoryError("UserDataRepo.RemoveAllByOwner", "commit", err)
			}
		}
	}()
	// стейтмент
	stmt, err := tx.PrepareContext(ctx, sqlUserDataRemoveAllByOwner)
	if err != nil {
		return apperrs.NewDalRepositoryError("UserDataRepo.RemoveAllByOwner", "remove sql statement", err)
	}
	defer stmt.Close()

	err = udrp.execStmt(ctx, stmt, userID)
	if err != nil {
		return apperrs.NewDalRepositoryError("UserDataRepo.RemoveAllByOwner", "update data", err)
	}

	return nil
}
