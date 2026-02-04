package model

import (
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	"github.com/google/uuid"
)

type UserKey struct {
	Username string
}

func NewUserKey(username string) *UserKey {
	return &UserKey{Username: username}
}

func NewEmptyUserKey() *UserKey {
	return NewUserKey("")
}

func (uk *UserKey) Validate() error {
	if uk.Username == "" {
		return apperrs.NewBllModelValidateError("UserKey.Username", "User must not be empty")
	}

	return nil
}

type User struct {
	ID           string
	Key          *UserKey
	PasswordHash string
	PrivateKey   string
	PublicKey    string
	Active       bool
	Deleted      bool
	Person       string
	EMail        string
	Data         []*UserData
	changed      bool
	ownedChanged bool
}

func NewUser(id string, username string, passwordHash string, privateKey string, publicKey string, active bool, person, eMail string) *User {
	return &User{
		ID:           id,
		Key:          NewUserKey(username),
		PasswordHash: passwordHash,
		PrivateKey:   privateKey,
		PublicKey:    publicKey,
		Active:       active,
		Data:         make([]*UserData, 0),
		Person:       person,
		EMail:        eMail,
		changed:      false,
		ownedChanged: false,
	}
}

func NewEmptyUser() *User {
	return &User{
		Key:          NewEmptyUserKey(),
		Data:         make([]*UserData, 0),
		changed:      false,
		ownedChanged: false,
	}
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) GetID() string {
	return u.ID
}

func (u *User) IsExists() bool {
	return u.ID != ""
}

func (u *User) ValidateCreate() error {
	if u.ID != "" {
		return apperrs.NewBllModelValidateError("User.ID", "must not be set")
	}
	if err := u.Key.Validate(); err != nil {
		return err
	}
	if err := u.validateAttrs(); err != nil {
		return err
	}

	return nil
}

func (u *User) ValidateChange() error {
	if u.ID == "" {
		return apperrs.NewBllModelValidateError("User.ID", "must not be set")
	}
	if err := u.Key.Validate(); err != nil {
		return err
	}
	if err := u.validateAttrs(); err != nil {
		return err
	}

	return nil
}

func (u *User) validateAttrs() error {
	if u.PasswordHash == "" {
		return apperrs.NewBllModelValidateError("User.PasswordHash", "must not be empty")
	}
	if u.PrivateKey == "" {
		return apperrs.NewBllModelValidateError("User.PrivateKey", "must not be empty")
	}
	if u.PublicKey == "" {
		return apperrs.NewBllModelValidateError("User.PublicKey", "must not be empty")
	}

	return nil
}

func (u *User) BeforeCreate() error {
	newID, err := uuid.NewRandom()
	if err != nil {
		return apperrs.NewBllModelError("user", "generate new id", err)
	}

	u.ID = newID.String()

	return nil
}

func (u *User) BeforeChange() error {
	return nil
}

func (u *User) IsDeleted() bool {
	return u.Deleted
}
