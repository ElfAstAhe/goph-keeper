package model

import (
	"time"

	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	"github.com/google/uuid"
)

type UserDataKey struct {
	Name     string
	DataKind string
}

func NewUserDataKey(name, dataKind string) *UserDataKey {
	return &UserDataKey{
		Name:     name,
		DataKind: dataKind,
	}
}

func NewEmptyUserDataKey() *UserDataKey {
	return NewUserDataKey("", "")
}

func (udk *UserDataKey) Validate() error {
	if udk.Name == "" {
		return apperrs.NewDalValidateError("UserDataKey.Name", "empty", nil)
	}
	if !IsDataKind(udk.DataKind) {
		return apperrs.NewDalValidateError("UserDataKey.DataKind", "mismatch", nil)
	}

	return nil
}

type UserData struct {
	UserDataShort
	TextData   string
	BinaryData []byte
}

func NewUserData(id string, name string, dataKind string, textData string, binaryData []byte, deleted bool) *UserData {
	return &UserData{
		UserDataShort: *NewUserDataShort(id, NewUserDataKey(name, dataKind), time.Now(), time.Now(), deleted),
		TextData:      textData,
		BinaryData:    binaryData,
	}
}

func NewEmptyUserData() *UserData {
	return &UserData{
		UserDataShort: *NewEmptyUserDataShort(),
		BinaryData:    make([]byte, 0),
	}
}

func (ud *UserData) TableName() string {
	return "user_data"
}

func (ud *UserData) GetID() string {
	return ud.ID
}

func (ud *UserData) IsExists() bool {
	return ud.ID != ""
}

func (ud *UserData) ValidateCreate() error {
	if ud.ID != "" {
		return apperrs.NewDalValidateError("UserData.ID", "not empty", nil)
	}
	if err := ud.Key.Validate(); err != nil {
		return err
	}

	return nil
}

func (ud *UserData) ValidateChange() error {
	if ud.ID == "" {
		return apperrs.NewDalValidateError("UserData.ID", "empty", nil)
	}
	if err := ud.Key.Validate(); err != nil {
		return err
	}

	return nil
}

func (ud *UserData) BeforeCreate() error {
	newID, err := uuid.NewRandom()
	if err != nil {
		return apperrs.NewDalCommonError("UserData.BeforeCreate", "generate new id", err)
	}

	ud.ID = newID.String()
	if ud.CreatedAt.IsZero() {
		ud.CreatedAt = time.Now()
	}
	ud.ModifiedAt = time.Now()

	return nil
}

func (ud *UserData) BeforeChange() error {
	ud.ModifiedAt = time.Now()

	return nil
}

func (ud *UserData) IsDeleted() bool {
	return ud.Deleted
}
