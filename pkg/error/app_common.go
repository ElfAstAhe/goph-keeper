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
	msg := "CMN: error"
	if e.message != "" {
		msg = fmt.Sprintf("%s with message [%s]", msg, e.message)
	}
	if e.err != nil {
		msg = fmt.Sprintf("%s: %v", msg, e.err)
	}

	return msg
}

func (e *AppCommonError) Unwrap() error {
	return e.err
}
