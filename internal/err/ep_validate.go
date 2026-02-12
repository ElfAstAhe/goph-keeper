package err

import (
	"fmt"
)

// EpValidateError — ошибка проверки данные
type EpValidateError struct {
	Dto   string // Какая dto (например, "User" или "UserData")
	Value string // Какое значение вызвало конфликт (например, "login 'admin'")
	Msg   string // Дополнительное сообщение
	Err   error  // Исходная ошибка из драйвера БД (опционально)
}

var ErrEpValidate *EpValidateError

func NewEpValidateError(dto, value string, Msg string, err error) *EpValidateError {
	return &EpValidateError{
		Dto:   dto,
		Value: value,
		Msg:   Msg,
		Err:   err,
	}
}

func (e *EpValidateError) Error() string {
	msg := fmt.Sprintf("EP: %s with value [%s] validation failed", e.Dto, e.Value)
	if e.Value == "" {
		msg = fmt.Sprintf("EP: %s validation failed", e.Dto)
	}
	if e.Msg != "" {
		msg = fmt.Sprintf("%s with message %s", msg, e.Msg)
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", msg, e.Err)
	}

	return msg
}

func (e *EpValidateError) Unwrap() error {
	return e.Err
}
