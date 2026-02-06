package err

import (
	"fmt"
)

// DalSoftDeletedError — сущность есть, но её удалили
type DalSoftDeletedError struct {
	Entity string // Какая сущность (например, "User" или "UserData")
	Key    string // Ключ сущности, которая была удалена
}

func NewDalSoftDeletedError(entity, key string) *DalSoftDeletedError {
	return &DalSoftDeletedError{
		Entity: entity,
		Key:    key,
	}
}

func (e *DalSoftDeletedError) Error() string {
	msg := fmt.Sprintf("DAL: %s with key [%s] soft deleted", e.Entity, e.Key)

	return msg
}
