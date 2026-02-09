package model

import (
	"time"
)

type UserDataShort struct {
	ID         string
	Key        *UserDataKey
	CreatedAt  time.Time
	ModifiedAt time.Time
	Deleted    bool
}

func NewUserDataShort(id string, key *UserDataKey, createdAt time.Time, modifiedAt time.Time, deleted bool) *UserDataShort {
	return &UserDataShort{
		ID:         id,
		Key:        key,
		CreatedAt:  createdAt,
		ModifiedAt: modifiedAt,
		Deleted:    deleted,
	}
}

func NewEmptyUserDataShort() *UserDataShort {
	return NewUserDataShort("", NewEmptyUserDataKey(), time.Now(), time.Now(), false)
}

func (uds *UserDataShort) TableName() string {
	return "user_data"
}

func (uds *UserDataShort) GetID() string {
	return uds.ID
}

func (uds *UserDataShort) IsExists() bool {
	return uds.ID != ""
}

func (uds *UserDataShort) IsDeleted() bool {
	return uds.Deleted
}
