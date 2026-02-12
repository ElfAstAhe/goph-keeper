package err

import (
	"fmt"
	"runtime"
)

// BllCommonError - обобщённая ошибка слоя BLL, приводит к InternalServerError
type BllCommonError struct {
	Msg  string
	Err  error
	File string
	Line int
}

func NewBllCommonError(msg string, err error) *BllCommonError {
	e := &BllCommonError{
		Msg: msg,
		Err: err,
	}
	// runtime.Caller(1) берет данные о том, КТО вызвал NewBllCommonError
	_, file, line, ok := runtime.Caller(1)
	if ok {
		e.File = file
		e.Line = line
	}

	return e
}

func (e *BllCommonError) Error() string {
	stack := ""
	if e.File != "" {
		// Формат [file.go:123] удобен для IDE (можно кликнуть в консоли)
		stack = fmt.Sprintf("[%s:%d] ", e.File, e.Line)
	}

	if e.Err != nil {
		return fmt.Sprintf("%sBLL: %s: %v", stack, e.Msg, e.Err)
	}

	return fmt.Sprintf("%sBLL: %s", stack, e.Msg)
}

func (e *BllCommonError) Unwrap() error {
	return e.Err
}
