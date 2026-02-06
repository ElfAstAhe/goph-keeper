package err

import "fmt"

// DalValidateError — ошибка уникальности данные
type DalValidateError struct {
	Entity string // Какая сущность (например, "User" или "UserData")
	Value  string // Какое значение вызвало конфликт (например, "login 'admin'")
	Err    error  // Исходная ошибка из драйвера БД (опционально)
}

func NewDalValidateError(entity, value string, err error) *DalValidateError {
	return &DalValidateError{
		Entity: entity,
		Value:  value,
		Err:    err,
	}
}

func (e *DalValidateError) Error() string {
	msg := fmt.Sprintf("DAL: %s with value [%s] validation failed", e.Entity, e.Value)
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", msg, e.Err)
	}

	return msg
}

func (e *DalValidateError) Unwrap() error {
	return e.Err
}
