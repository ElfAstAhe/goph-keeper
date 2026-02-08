package err

import (
	"fmt"
)

// BllValidateError — ошибка проверки данные
type BllValidateError struct {
	Entity string // Какая сущность (например, "User" или "UserData")
	Value  string // Какое значение вызвало конфликт (например, "login 'admin'")
	Msg    string // Дополнительное сообщение
	Err    error  // Исходная ошибка из драйвера БД (опционально)
}

var ErrBllValidate *BllValidateError

func NewBllValidateError(entity, value string, Msg string, err error) *BllValidateError {
	return &BllValidateError{
		Entity: entity,
		Value:  value,
		Msg:    Msg,
		Err:    err,
	}
}

func (e *BllValidateError) Error() string {
	msg := fmt.Sprintf("BLL: %s with value [%s] validation failed", e.Entity, e.Value)
	if e.Value == "" {
		msg = fmt.Sprintf("BLL: %s validation failed", e.Entity)
	}
	if e.Msg != "" {
		msg = fmt.Sprintf("%s with message %s", msg, e.Msg)
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", msg, e.Err)
	}

	return msg
}

func (e *BllValidateError) Unwrap() error {
	return e.Err
}
