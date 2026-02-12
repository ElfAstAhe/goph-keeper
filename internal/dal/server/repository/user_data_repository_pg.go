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
and user_id = $2
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
and user_id = $2
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

	sqlUserDataListAllByOwner string = `
select
    id,
    name,
    kind,
    created_at,
    modified_at,
    deleted
from
    user_data
where
    user_id = $1
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

func (udrp *UserDataRepositoryPg) Get(ctx context.Context, userID, ID string) (*model.UserData, error) {
	res, err := udrp.internalGetSingle(ctx, sqlUserDataGet, ID, userID)
	if err != nil {
		return nil, err
	}

	return udrp.afterGet(res)
}

func (udrp *UserDataRepositoryPg) GetByKey(ctx context.Context, userID string, key *model.UserDataKey) (*model.UserData, error) {
	res, err := udrp.internalGetSingle(ctx, sqlUserDataGetByKey, userID, key.Name, key.DataKind)
	if err != nil {
		return nil, apperrs.NewDalCommonError("UserDataRepo.GetByKey", "get by key", err)
	}

	return udrp.afterGet(res)
}

func (udrp *UserDataRepositoryPg) internalGetSingle(ctx context.Context, sqlReq string, params ...any) (*model.UserData, error) {
	row := udrp.db.GetDB().QueryRowContext(ctx, sqlReq, params...)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, apperrs.NewDalCommonError("UserDataRepo.internalGetSingle", "fetch row", row.Err())
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
			return nil, apperrs.NewDalNotFoundError("UserDataRepo.internalGetSingle", "user data not found", err)
		}

		return nil, apperrs.NewDalCommonError("UserDataRepo.internalGetSingle", "scan row", err)
	}

	return entity, nil
}

func (udrp *UserDataRepositoryPg) afterGet(res *model.UserData) (*model.UserData, error) {
	var err error = nil
	// расшифровываем данные
	res.TextData = udrp.dataCipherHelper.DecryptString(res.TextData)
	res.BinaryData = udrp.dataCipherHelper.DecryptBinary(res.BinaryData)

	if res.Deleted {
		err = apperrs.NewDalSoftDeletedError("UserData", fmt.Sprintf("ID: [%s], Key: [%v]", res.ID, res.Key))
	}

	return res, err
}

func (udrp *UserDataRepositoryPg) afterGetMulti(res *model.UserDataShort) (*model.UserDataShort, error) {
	var err error = nil
	if res.Deleted {
		err = apperrs.NewDalSoftDeletedError("UserDataShort", fmt.Sprintf("ID: [%s]", res.ID))
	}

	return res, err
}

func (udrp *UserDataRepositoryPg) Create(ctx context.Context, userID string, userData *model.UserData) (res *model.UserData, err error) {
	// валидируем
	if err = udrp.validateCreate(userID, userData); err != nil {
		return nil, apperrs.NewDalValidateError("UserData", "validate create", err)
	}
	// подготавливаем
	if err = udrp.beforeCreate(userData); err != nil {
		return nil, apperrs.NewDalCommonError("UserDataRepo.Create", "before create", err)
	}

	// сохраняем
	err = udrp.db.GetHelper().RunInTx(ctx, udrp.db.GetDB(), func(tx *sql.Tx) error {
		// стейтмент
		stmt, err := tx.PrepareContext(ctx, sqlUserDataCreate)
		if err != nil {
			return apperrs.NewDalCommonError("UserDataRepo.Create", "prepare stmt", err)
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx,
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
			if udrp.db.GetHelper().IsUniqueViolation(err) {
				return apperrs.NewDalAlreadyExistsError("UserData", fmt.Sprintf("userID [%s], key [%v]", userID, userData.Key), err)
			}

			return apperrs.NewDalCommonError("UserDataRepo.Create", "exec stmt", err)
		}

		return nil
	})
	if err != nil {
		return nil, apperrs.NewDalCommonError("UserDataRepo.Create", "run in transaction", err)
	}

	return udrp.afterGet(userData)
}

func (udrp *UserDataRepositoryPg) validateCreate(userID string, userData *model.UserData) error {
	if userData == nil {
		return errs.NewAppInvalidArgumentError("userData", "nil user data")
	}
	if userID == "" {
		return errs.NewAppInvalidArgumentError("userData", "empty user id")
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
		return nil, apperrs.NewDalValidateError("UserData", "validate change", err)
	}
	// подготавливаем
	if err = udrp.beforeChange(userData); err != nil {
		return nil, apperrs.NewDalCommonError("UserDataRepo.Change", "before change", err)
	}

	// сохраняем
	err = udrp.db.GetHelper().RunInTx(ctx, udrp.db.GetDB(), func(tx *sql.Tx) error {
		// стейтмент
		stmt, err := tx.PrepareContext(ctx, sqlUserDataChange)
		if err != nil {
			return apperrs.NewDalCommonError("UserDataRepo.Change", "prepare stmt", err)
		}
		defer stmt.Close()

		err = udrp.db.GetHelper().ExecStmt(ctx, stmt, func(repErr error) (string, string, error) {
			return "UserData", userData.GetID(), repErr
		},
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
			return apperrs.NewDalCommonError("UserDataRepo.Change", "exec stmt", err)
		}

		return nil
	})
	if err != nil {
		return nil, apperrs.NewDalCommonError("UserDataRepo.Change", "run in transaction", err)
	}

	return udrp.afterGet(userData)
}

func (udrp *UserDataRepositoryPg) validateChange(userID string, userData *model.UserData) error {
	if userData == nil {
		return apperrs.NewDalCommonError("UserDataRepo.validateChange", "nil user data instance", nil)
	}
	if userID == "" {
		return apperrs.NewDalCommonError("UserDataRepo.validateChange", "empty user ID", nil)
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

func (udrp *UserDataRepositoryPg) Remove(ctx context.Context, userID, id string) error {
	// сохраняем
	err := udrp.db.GetHelper().RunInTx(ctx, udrp.db.GetDB(), func(tx *sql.Tx) error {
		// стейтмент
		stmt, err := tx.PrepareContext(ctx, sqlUserDataRemove)
		if err != nil {
			return apperrs.NewDalCommonError("UserDataRepo.Remove", "remove sql statement", err)
		}
		defer stmt.Close()

		err = udrp.db.GetHelper().ExecStmt(ctx, stmt, func(repErr error) (string, string, error) {
			return "UserData", id, repErr
		}, id, userID)
		if err != nil {
			return apperrs.NewDalCommonError("UserDataRepo.Remove", "update data", err)
		}

		return nil
	})
	if err != nil {
		return apperrs.NewDalCommonError("UserDataRepo.Remove", "run transaction", err)
	}

	return nil
}

func (udrp *UserDataRepositoryPg) ListAllByOwner(ctx context.Context, userID string) ([]*model.UserDataShort, error) {
	return udrp.internalGetMulti(ctx, sqlUserDataListAllByOwner, userID)
}

func (udrp *UserDataRepositoryPg) internalGetMulti(ctx context.Context, sqlReq string, params ...any) ([]*model.UserDataShort, error) {
	rows, err := udrp.db.GetDB().QueryContext(ctx, sqlReq, params...)
	if err != nil {
		return nil, apperrs.NewDalCommonError("UserDataRepo.ListAllByOwner", "query", err)
	}
	defer rows.Close()

	res := make([]*model.UserDataShort, 0)
	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return nil, apperrs.NewDalCommonError("UserDataRepo.ListAllByOwner", "check context", err)
		}

		entity := model.NewEmptyUserDataShort()

		err := rows.Scan(&entity.ID,
			&entity.Key.Name,
			&entity.Key.DataKind,
			&entity.CreatedAt,
			&entity.ModifiedAt,
			&entity.Deleted,
		)
		if err != nil {
			return nil, apperrs.NewDalCommonError("UserDataRepo.ListAllByOwner", "scan rows", err)
		}

		entity, err = udrp.afterGetMulti(entity)
		if err != nil {
			if errors.As(err, &apperrs.ErrDalSoftDeleted) {
				continue
			}

			return nil, apperrs.NewDalCommonError("UserDataRepo.ListAllByOwner", "post scan processing", err)
		}

		res = append(res, entity)
	}
	if rows.Err() != nil {
		return nil, apperrs.NewDalCommonError("UserDataRepo.ListAllByOwner", "after scan", rows.Err())
	}

	return res, nil
}

func (udrp *UserDataRepositoryPg) RemoveAllByOwner(ctx context.Context, userID string) error {
	// сохраняем
	err := udrp.db.GetHelper().RunInTx(ctx, udrp.db.GetDB(), func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, sqlUserDataRemoveAllByOwner)
		if err != nil {
			return apperrs.NewDalCommonError("UserDataRepo.RemoveAllByOwner", "remove sql statement", err)
		}
		defer stmt.Close()

		err = udrp.db.GetHelper().ExecStmt(ctx, stmt, func(repErr error) (string, string, error) {
			return "UserData", userID, repErr
		}, userID)
		if err != nil {
			if !errors.As(err, &apperrs.ErrDalNotFound) {
				return apperrs.NewDalCommonError("UserDataRepo.RemoveAllByOwner", "update data", err)
			}
		}

		return nil
	})
	if err != nil {
		return apperrs.NewDalCommonError("UserDataRepo.RemoveAllByOwner", "run transaction", err)
	}

	return nil
}
