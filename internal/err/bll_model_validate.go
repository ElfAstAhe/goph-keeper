package err

import "fmt"

type BllModelValidateError struct {
	name string
	msg  string
}

var ErrBllModelValidate *BllModelValidateError

func NewBllModelValidateError(name string, msg string) *BllModelValidateError {
	return &BllModelValidateError{
		name: name,
		msg:  msg,
	}
}

func (mv *BllModelValidateError) Error() string {
	return fmt.Sprintf("BLL: model [%s] validation error: [%s]", mv.name, mv.msg)
}
