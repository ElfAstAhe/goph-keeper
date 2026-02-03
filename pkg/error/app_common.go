package error

import (
	"fmt"
)

type AppCommonError struct {
	message string
	err     error
}

var ErrAppCommon *AppCommonError

func NewAppCommonError(msg string, err error) *UtlCipherError {
	return &UtlCipherError{
		message: msg,
		err:     err,
	}
}

func (e *AppCommonError) Error() string {
	return fmt.Sprintf("cipher error with message [%s], with error [%s]", e.message, e.err)
}

func (e *AppCommonError) Unwrap() error {
	return e.err
}
