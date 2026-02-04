package model

import (
	"time"

	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	"github.com/google/uuid"
)

type UserDataKey struct {
	UserID   string
	Name     string
	DataKind string
}

func NewUserDataKey(userID, name, dataKind string) *UserDataKey {
	return &UserDataKey{
		UserID:   userID,
		Name:     name,
		DataKind: dataKind,
	}
}

func NewEmptyUserDataKey() *UserDataKey {
	return NewUserDataKey("", "", "")
}

func (udk *UserDataKey) Validate() error {
	if udk.UserID == "" {
		return apperrs.NewBllModelValidateError("UserDataKey.UserID", "must be set")
	}
	if udk.Name == "" {
		return apperrs.NewBllModelValidateError("UserDataKey.Name", "must be set")
	}
	if !IsDataKind(udk.DataKind) {
		return apperrs.NewBllModelValidateError("UserDataKey.DataKind", "must be set")
	}

	return nil
}

type UserData struct {
	ID         string
	Key        *UserDataKey
	TextData   string
	BinaryData []byte
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	Deleted    bool
}

func NewUserData(id string, userID string, name string, dataKind string, textData string, binaryData []byte, deleted bool) *UserData {
	currentTime := time.Now()
	return &UserData{
		ID:         id,
		Key:        NewUserDataKey(userID, name, dataKind),
		TextData:   textData,
		BinaryData: binaryData,
		CreatedAt:  &currentTime,
		UpdatedAt:  &currentTime,
		Deleted:    deleted,
	}
}

func NewEmptyOwnedUserData(userID string) *UserData {
	return &UserData{
		Key: NewUserDataKey(userID, "", ""),
	}
}

func NewEmptyUserData() *UserData {
	return NewEmptyOwnedUserData("")
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
		return apperrs.NewBllModelValidateError("UserData.ID", "must not be set")
	}
	if err := ud.Key.Validate(); err != nil {
		return err
	}

	return nil
}

func (ud *UserData) ValidateChange() error {
	if ud.ID == "" {
		return apperrs.NewBllModelValidateError("UserData.ID", "must not be set")
	}
	if err := ud.Key.Validate(); err != nil {
		return err
	}

	return nil
}

func (ud *UserData) BeforeCreate() error {
	newID, err := uuid.NewRandom()
	if err != nil {
		return apperrs.NewBllModelError("user", "generate new id", err)
	}

	ud.ID = newID.String()

	return nil
}

func (ud *UserData) BeforeChange() error {
	return nil
}

func (ud *UserData) IsDeleted() bool {
	return ud.Deleted
}
