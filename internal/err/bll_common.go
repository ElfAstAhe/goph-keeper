package err

import "fmt"

type BllCommonError struct {
	msg string
	err error
}

var ErrBllCommon *BllCommonError

func NewBllCommonError(msg string, err error) *BllCommonError {
	return &BllCommonError{msg: msg, err: err}
}

func (e *BllCommonError) Error() string {
	return fmt.Sprintf("BLL: error performing operation with message [%s] with error [%v]", e.msg, e.err.Error())
}

func (e *BllCommonError) Unwrap() error {
	return e.err
}
