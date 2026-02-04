package err

import (
	"fmt"
)

type BllModelError struct {
	model string
	msg   string
	err   error
}

var ErrBllModel *BllModelError

func NewBllModelError(model string, msg string, err error) *BllModelError {
	return &BllModelError{
		model: model,
		msg:   msg,
		err:   err,
	}
}

func (err *BllModelError) Error() string {
	return fmt.Sprintf("BLL: model [%s] with message [%s] with error [%v]", err.model, err.msg, err.err)
}

func (err *BllModelError) Unwrap() error {
	return err.err
}
