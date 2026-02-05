package err

import (
	"fmt"
)

type DalCommonError struct {
	Op  string // Метод: "UserRepo.Create"
	Msg string // Суть: "scan row" или "execute query"
	Err error  // Исходная ошибка из database/sql или драйвера
}

func NewDalCommonError(op, msg string, err error) *DalCommonError {
	return &DalCommonError{Op: op, Msg: msg, Err: err}
}

func (e *DalCommonError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("DAL: [%s] %s: %v", e.Op, e.Msg, e.Err)
	}

	return fmt.Sprintf("DAL: [%s] %s", e.Op, e.Msg)
}

func (e *DalCommonError) Unwrap() error {
	return e.Err
}
