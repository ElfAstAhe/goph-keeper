package err

import "fmt"

// DalNotFoundError — отсутствие сущности
type DalNotFoundError struct {
    Entity string // Какая сущность (например, "User" или "UserData")
    Value  string // Какое значение вызвало конфликт (например, "login 'admin'")
    Err    error  // Исходная ошибка из драйвера БД (опционально)
}

func NewDalNot                                                                                                 FoundError(entity, value string, err error) *DalNotFoundError {
    return &DalNotFoundError{
        Entity: entity,
        Value:  value,
        Err:    err,
    }
}

func (e *DalNotFoundError) Error() string {
    msg := fmt.Sprintf("DAL: %s with value [%s] already exists", e.Entity, e.Value)
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", msg, e.Err)
    }

    return msg
}

func (e *DalNotFoundError) Unwrap() error {
    return e.Err
}
