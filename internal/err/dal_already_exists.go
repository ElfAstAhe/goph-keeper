package err

import "fmt"

// DalAlreadyExistsError — ошибка уникальности данные
type DalAlreadyExistsError struct {
	Entity string // Какая сущность (например, "User" или "UserData")
	Value  string // Какое значение вызвало конфликт (например, "login 'admin'")
	Err    error  // Исходная ошибка из драйвера БД (опционально)
}

func NewDalAlreadyExistsError(entity, value string, err error) *DalAlreadyExistsError {
	return &DalAlreadyExistsError{
		Entity: entity,
		Value:  value,
		Err:    err,
	}
}

func (e *DalAlreadyExistsError) Error() string {
	msg := fmt.Sprintf("DAL: %s with value [%s] already exists", e.Entity, e.Value)
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", msg, e.Err)
	}

	return msg
}

func (e *DalAlreadyExistsError) Unwrap() error {
	return e.Err
}
